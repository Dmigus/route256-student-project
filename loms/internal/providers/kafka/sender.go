// Package kafka содержит Sender для записи в топик кафки
package kafka

import (
	"context"
	"github.com/IBM/sarama"
	"route256.ozon.ru/project/loms/internal/models"
)

const atHeaderKey = "At"

// Sender это провайдер, который умеет отправлять сообщения в кафку
type Sender struct {
	topic    string
	producer sarama.SyncProducer
}

// NewSender создайт новый Sender
func NewSender(brokers []string, topic string) (*Sender, error) {
	syncProducer, err := sarama.NewSyncProducer(brokers, getConfig())
	if err != nil {
		return nil, err
	}
	return &Sender{
		producer: syncProducer,
		topic:    topic,
	}, nil
}

// SendMessages синхронно отправляет сообщения events в брокер
func (p *Sender) SendMessages(_ context.Context, messages []models.EventMessage) error {
	saramaMessages := make([]*sarama.ProducerMessage, 0, len(messages))
	for _, ev := range messages {
		message := p.modelMessageToSarama(ev)
		saramaMessages = append(saramaMessages, message)
	}
	return p.producer.SendMessages(saramaMessages)
}

func (p *Sender) modelMessageToSarama(ev models.EventMessage) *sarama.ProducerMessage {
	return &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.ByteEncoder(ev.PartitionKey),
		Value: sarama.ByteEncoder(ev.Payload),
	}
}

func getConfig() *sarama.Config {
	c := sarama.NewConfig()
	c.Producer.Partitioner = sarama.NewHashPartitioner
	c.Producer.Return.Successes = true
	c.Producer.Return.Errors = true
	// at least once
	c.Producer.RequiredAcks = sarama.WaitForAll
	c.Net.MaxOpenRequests = 1
	return c
}
