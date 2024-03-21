package models

import "errors"

type OrderStatus int

var ErrWrongOrderStatus = errors.New("order status is wrong")

const (
	New OrderStatus = iota + 1
	AwaitingPayment
	Failed
	Payed
	Cancelled
)

type Order struct {
	id              int64
	UserId          int64
	Status          OrderStatus
	Items           []OrderItem
	IsItemsReserved bool
}

func NewOrder(userId, orderId int64) *Order {
	return &Order{
		UserId:          userId,
		id:              orderId,
		Status:          New,
		IsItemsReserved: false,
	}
}

func (o *Order) Id() int64 {
	return o.id
}
