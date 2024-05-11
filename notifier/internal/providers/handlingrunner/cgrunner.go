// Package handlingrunner содержит функциональность, которая позволяет запустить обработку событий, получаемых из некоторого источника
package handlingrunner

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/IBM/sarama"
	"github.com/dnwe/otelsarama"
	"go.uber.org/zap"
	"route256.ozon.ru/project/notifier/internal/service"
)

// KafkaConsumerGroupRunner это структура, которая умеет запускать обработку событий, получаемых из кафки
type KafkaConsumerGroupRunner struct {
	configVersion atomic.Uint64
	cg            *reconfigurableConsumerGroup
	mu            sync.Mutex
	logger        *zap.Logger
	topic         string
}

// NewKafkaConsumerGroupRunner возращает новый KafkaConsumerGroupRunner, сконфигурированный на брокеры brokers и топик topic
func NewKafkaConsumerGroupRunner(config ConsumerGroupRunnerConfig) (*KafkaConsumerGroupRunner, error) {
	runner := &KafkaConsumerGroupRunner{cg: newReconfigurableConsumerGroup()}
	err := runner.initFromConfig(config)
	if err != nil {
		return nil, err
	}
	return runner, nil
}

func (k *KafkaConsumerGroupRunner) initFromConfig(config ConsumerGroupRunnerConfig) error {
	err := k.cg.Init(config.Brokers, getConfig())
	if err != nil {
		return err
	}
	k.mu.Lock()
	defer k.mu.Unlock()
	k.logger = config.Logger
	k.topic = config.Topic
	k.configVersion.Add(1)
	return nil
}

// Update это блокируюший вызов, который обновляет KafkaConsumerGroupRunner в соответствии переданным config.
// Если не удалось обновить с таким конфигом, предыдущая версия конфигурации остаётся работать
func (k *KafkaConsumerGroupRunner) Update(config ConsumerGroupRunnerConfig) error {
	return k.initFromConfig(config)
}

// Run обрабатывает поступающие сообщения переданным хандлером в рамках группы. Блокирующий.
func (k *KafkaConsumerGroupRunner) Run(ctx context.Context, handler service.EventHandler) (err error) {
	defer k.cg.Close()
	k.mu.Lock()
	saramaHandler := newConsumerGroupHandler(handler, k.logger)
	k.mu.Unlock()
	tracedHandler := otelsarama.WrapConsumerGroupHandler(saramaHandler)
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			k.consumeCycle(ctx, &tracedHandler)
		}
	}
}

func (k *KafkaConsumerGroupRunner) consumeCycle(ctx context.Context, handler *sarama.ConsumerGroupHandler) {
	k.mu.Lock()
	configVer := k.configVersion.Load()
	topic := k.topic
	k.mu.Unlock()
	cg := k.cg.GetInitializedCG(ctx)
	// проверяем, что полученный cg всё ещё актуален
	if cg == nil || configVer != k.configVersion.Load() {
		return
	}
	err := cg.Consume(ctx, []string{topic}, *handler)
	if err != nil {
		k.mu.Lock()
		logger := k.logger
		k.mu.Unlock()
		logger.Error("error in consumer group session", zap.Error(err))
	}
}

// ConsumerGroupRunnerConfig это конфигурация для KafkaConsumerGroupRunner
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
