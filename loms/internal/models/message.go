package models

import "context"

// EventMessage это структура, представляющая событие изменения статуса заказа
type EventMessage struct {
	PartitionKey []byte
	Payload      []byte
	TraceContext context.Context
}
