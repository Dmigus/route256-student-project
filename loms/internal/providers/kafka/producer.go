// Package kafka содержит Producer для записи в топик кафки
package kafka

import (
	"context"
	"strconv"
	"time"

	"github.com/IBM/sarama"
	"route256.ozon.ru/project/loms/internal/models"
)

// Producer это провайдер, который умеет отправлять события изменения статуса заказа в брокер кафки
type Producer struct {
	topic    string
	producer sarama.SyncProducer
}

// NewProducer создайт новый Producer
func NewProducer(brokers []string, topic string) (*Producer, error) {
	config := PrepareConfig(
		WithProducerPartitioner(sarama.NewHashPartitioner),
		WithRequiredAcks(sarama.WaitForAll),
		WithMaxOpenRequests(1),
		WithMaxRetries(5),
		WithRetryBackoff(10*time.Millisecond),
	)
	syncProducer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}
	return &Producer{
		producer: syncProducer,
		topic:    topic,
	}, nil
}

// PushOrderStatusChangedEvents синхронно отправляет сообщения events в брокер
func (p *Producer) PushOrderStatusChangedEvents(_ context.Context, events []models.OrderStatusChangedEvent) error {
	messages := make([]*sarama.ProducerMessage, 0, len(events))
	for _, ev := range events {
		message := p.orderStatusChangeEvToMessage(ev)
		messages = append(messages, message)
	}
	return p.producer.SendMessages(messages)
}

func (p *Producer) orderStatusChangeEvToMessage(ev models.OrderStatusChangedEvent) *sarama.ProducerMessage {
	evTimeStampHeader := sarama.RecordHeader{Key: []byte("At"), Value: []byte(ev.At.String())}
	return &sarama.ProducerMessage{
		Topic:   p.topic,
		Key:     sarama.StringEncoder(strconv.FormatInt(ev.OrderID, 10)),
		Value:   sarama.ByteEncoder(ev.Message),
		Headers: []sarama.RecordHeader{evTimeStampHeader},
	}
}
