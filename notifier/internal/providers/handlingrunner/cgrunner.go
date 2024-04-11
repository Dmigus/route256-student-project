package handlingrunner

import (
	"context"
	"github.com/IBM/sarama"
	"route256.ozon.ru/project/notifier/internal/service"
)

const groupName = "notifier-group"

type KafkaConsumerGroupRunner struct {
	sarama.ConsumerGroup
	topic string
}

func NewKafkaConsumerGroupRunner(brokers []string, topic string) (*KafkaConsumerGroupRunner, error) {
	cg, err := sarama.NewConsumerGroup(brokers, groupName, getConfig())
	if err != nil {
		return nil, err
	}
	return &KafkaConsumerGroupRunner{
		ConsumerGroup: cg,
		topic:         topic,
	}, nil
}

// Run обрабатывает поступающие сообщения переданным хандлером. Блокирующий.
func (k *KafkaConsumerGroupRunner) Run(ctx context.Context, handler service.EventHandler) error {
	saramaHandler := newConsumerGroupHandler(handler)
	for {
		// Игнорируем ошибки
		_ = k.ConsumerGroup.Consume(ctx, []string{k.topic}, saramaHandler)
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

func getConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Offsets.AutoCommit.Enable = true
	return config
}
