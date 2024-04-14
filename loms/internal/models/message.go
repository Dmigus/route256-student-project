package models

// EventMessage это структура, представляющая событие изменения статуса заказа
type EventMessage struct {
	PartitionKey []byte
	Payload      []byte
}
