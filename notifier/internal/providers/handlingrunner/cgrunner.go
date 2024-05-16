// Package handlingrunner содержит функциональность, которая позволяет запустить обработку событий, получаемых из некоторого источника
package handlingrunner

import (
	"context"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"route256.ozon.ru/project/notifier/internal/service"
	"sync"
)

const groupName = "notifier-group"

// KafkaConsumerGroupRunner это структура, которая умеет запускать обработку событий, получаемых из кафки
type KafkaConsumerGroupRunner struct {
	mu     sync.Mutex
	config ConsumerGroupRunnerConfig
	cg     *underCounterCG
}

// NewKafkaConsumerGroupRunner возращает новый KafkaConsumerGroupRunner, сконфигурированный на брокеры brokers и топик topic
func NewKafkaConsumerGroupRunner(config ConsumerGroupRunnerConfig) (*KafkaConsumerGroupRunner, error) {
	return &KafkaConsumerGroupRunner{config: config}, nil
}

// init предназначен для инициализации коннекта
func (k *KafkaConsumerGroupRunner) init() error {
	k.mu.Lock()
	defer k.mu.Unlock()
	cg, err := sarama.NewConsumerGroup(k.config.Brokers, groupName, getConfig())
	if err != nil {
		return err
	}
	if k.cg != nil {
		k.cg.CloseWhenIsNotUsed()
	}
	k.cg = newUnderCounterCG(cg)
	return nil
}

// UpdateConfig это блокируюший вызов, который обновляет KafkaConsumerGroupRunner в соответствии переданным config.
// Если не удалось обновить с таким конфигом, предыдущая версия конфигурации остаётся работать
func (k *KafkaConsumerGroupRunner) UpdateConfig(config ConsumerGroupRunnerConfig) error {
	k.mu.Lock()
	k.config = config
	k.mu.Unlock()
	return k.init()
}

// Run обрабатывает поступающие сообщения переданным хандлером в рамках группы. Блокирующий.
func (k *KafkaConsumerGroupRunner) Run(ctx context.Context, handler service.EventHandler) (err error) {
	err = k.init()
	if err != nil {
		return err
	}
	defer k.cg.CloseWhenIsNotUsed()
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			k.consumeCycle(ctx, handler)
		}
	}
}

func (k *KafkaConsumerGroupRunner) consumeCycle(ctx context.Context, handler service.EventHandler) {
	k.mu.Lock()
	saramaHandler := newConsumerGroupHandler(handler, k.config.Logger)
	topic := k.config.Topic
	cgWrapper := k.cg
	cg := cgWrapper.GetForUsage()
	defer cgWrapper.Done()
	k.mu.Unlock()
	err := cg.Consume(ctx, []string{topic}, saramaHandler)
	if err != nil {
		k.mu.Lock()
		logger := k.config.Logger
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
