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
	Status OrderStatus
	Items  []OrderItem
}

func NewOrder(id int64) *Order {
	return &Order{
		id:     id,
		Status: New,
	}
}

func (o *Order) Id() int64 {
	return o.id
}
