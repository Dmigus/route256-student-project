// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package stocks

type ItemUnit struct {
	SkuID    int64
	Total    int32
	Reserved int32
}

type MessageOutbox struct {
	ID           int64
	PartitionKey []byte
	Payload      []byte
}

type Order struct {
	ID               int64
	UserID           int64
	Status           string
	AreItemsReserved bool
}

type OrderItem struct {
	OrderID int64
	SkuID   int64
	Count   int32
}
