package models

type OrderStatus int

const (
	New OrderStatus = iota + 1
	AwaitingPayment
	Failed
	Payed
	Cancelled
)

type Order struct {
	id     int64
	UserId int64
	Status OrderStatus
	Items  []OrderItem
}

func NewOrder(userId, orderId int64) *Order {
	return &Order{
		UserId: userId,
		id:     orderId,
		Status: New,
	}
}

func (o *Order) Id() int64 {
	return o.id
}
