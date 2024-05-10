// Package handlingrunner содержит функциональность, которая позволяет запустить обработку событий, получаемых из некоторого источника
package handlingrunner

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/dnwe/otelsarama"
	"go.uber.org/zap"
	"route256.ozon.ru/project/notifier/internal/service"
	"sync"
)

const groupName = "notifier-group"

// KafkaConsumerGroupRunner это структура, которая умеет запускать обработку событий, получаемых из кафки
type KafkaConsumerGroupRunner struct {
	mu            sync.Mutex
	cg            *sarama.ConsumerGroup
	cgInitialized chan struct{}
	logger        *zap.Logger
	handler       *sarama.ConsumerGroupHandler
	topic         string
}

// NewKafkaConsumerGroupRunner возращает новый KafkaConsumerGroupRunner, сконфигурированный на брокеры brokers и топик topic
func NewKafkaConsumerGroupRunner(config ConsumerGroupRunnerConfig) (*KafkaConsumerGroupRunner, error) {
	runner := &KafkaConsumerGroupRunner{cgInitialized: make(chan struct{})}
	err := runner.initFromConfig(config)
	if err != nil {
		return nil, err
	}
	return runner, nil
}

func (k *KafkaConsumerGroupRunner) initFromConfig(config ConsumerGroupRunnerConfig) error {
	cg, err := sarama.NewConsumerGroup(config.Brokers, groupName, getConfig())
	if err != nil {
		return err
	}
	k.mu.Lock()
	defer k.mu.Unlock()
	k.logger = config.Logger
	k.topic = config.Topic
	k.cg = &cg
	close(k.cgInitialized)
	return nil
}

func (k *KafkaConsumerGroupRunner) Update(config ConsumerGroupRunnerConfig) error {
	var err error
	k.mu.Lock()
	if k.isCGRunningLocked() {
		err = (*k.cg).Close()
		k.cgInitialized = make(chan struct{})
	}
	k.mu.Unlock()
	if err != nil {
		return err
	}
	return k.initFromConfig(config)
}

func (k *KafkaConsumerGroupRunner) isCGRunningLocked() bool {
	select {
	case <-k.cgInitialized:
		return true
	default:
		return false
	}
}

// Run обрабатывает поступающие сообщения переданным хандлером в рамках группы. Блокирующий.
func (k *KafkaConsumerGroupRunner) Run(ctx context.Context, handler service.EventHandler) (err error) {
	defer func() {
		k.mu.Lock()
		cg := *k.cg
		k.mu.Unlock()
		err = cg.Close()
	}()
	k.mu.Lock()
	saramaHandler := newConsumerGroupHandler(handler, k.logger)
	tracedHandler := otelsarama.WrapConsumerGroupHandler(saramaHandler)
	k.mu.Unlock()
	for {
		k.mu.Lock()
		cgInit := k.cgInitialized
		k.mu.Unlock()
		select {
		case <-ctx.Done():
			return nil
		default:
			select {
			case <-cgInit:
				k.oneIteration(ctx, &tracedHandler)
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func (k *KafkaConsumerGroupRunner) oneIteration(ctx context.Context, handler *sarama.ConsumerGroupHandler) {
	k.mu.Lock()
	cg := k.cg
	topic := k.topic
	k.mu.Unlock()
	err := (*cg).Consume(ctx, []string{topic}, *handler)
	if err != nil {
		k.mu.Lock()
		logger := k.logger
		k.mu.Unlock()
		logger.Error("error in consumer group session", zap.Error(err))
	}
}

type ConsumerGroupRunnerConfig struct {
	Brokers []string
	Topic   string
	Logger  *zap.Logger
}

func getConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Offsets.AutoCommit.Enable = true
	return config
}
