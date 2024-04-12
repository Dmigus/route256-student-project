// Package handlingrunner содержит функциональность, которая позволяет запустить обработку событий, получаемых из некоторого источника
package handlingrunner

import (
	"context"

	"github.com/IBM/sarama"
	"route256.ozon.ru/project/notifier/internal/service"
)

const groupName = "notifier-group"

// KafkaConsumerGroupRunner это структура, которая умеет запускать обработку событий, получаемых из кафки
type KafkaConsumerGroupRunner struct {
	brokers []string
	topic   string
}

// NewKafkaConsumerGroupRunner возращает новый KafkaConsumerGroupRunner, сконфигурированный на брокеры brokers и топик topic
func NewKafkaConsumerGroupRunner(brokers []string, topic string) *KafkaConsumerGroupRunner {
	return &KafkaConsumerGroupRunner{
		brokers: brokers,
		topic:   topic,
	}
}

// Run обрабатывает поступающие сообщения переданным хандлером в рамках группы. Блокирующий.
func (k *KafkaConsumerGroupRunner) Run(ctx context.Context, handler service.EventHandler) (err error) {
	cg, err := sarama.NewConsumerGroup(k.brokers, groupName, getConfig())
	if err != nil {
		return err
	}
	defer func() {
		err = cg.Close()
	}()
	saramaHandler := newConsumerGroupHandler(handler)
	for {
		// Непонятно куда логировать ошибки
		_ = cg.Consume(ctx, []string{k.topic}, saramaHandler)
		if ctx.Err() != nil {
			return nil
		}
	}
}

func getConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Offsets.AutoCommit.Enable = true
	return config
}
