package models

import "time"

// OrderStatusChangedEvent это структура, представляющая событие изменения статуса заказа
type OrderStatusChangedEvent struct {
	OrderID int64
	Message string
	At      time.Time
}
