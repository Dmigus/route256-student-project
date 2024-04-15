// Package models содержит сущности, которыми оперируют различные слои
package models

// EventMessage это структура, представляющая событие изменения статуса заказа
type EventMessage struct {
	PartitionKey []byte
	Payload      []byte
}
