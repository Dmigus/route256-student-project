package models

import "time"

type OrderStatusChangedEvent struct {
	OrderID int64
	Message string
	At      time.Time
}
