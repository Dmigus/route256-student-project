package modifier

import (
	"context"
	"route256.ozon.ru/project/loms/internal/models"
)

type Orders struct {
	queries *Queries
}

func NewOrders(tx DBTX) *Orders {
	return &Orders{queries: New(tx)}
}

// Create создаёт заказ для юзера userID и товарами items в репозитории и возращает его
func (po *Orders) Create(ctx context.Context, userID int64, items []models.OrderItem) (*models.Order, error) {
	params := createOrderParams{UserID: userID, Status: OrderStatusNew, AreItemsReserved: false}
	orderId, err := po.queries.createOrder(ctx, params)
	if err != nil {
		return nil, err
	}
	order := models.NewOrder(userID, orderId)
	itemsParams := insertItemParamsFrom(orderId, items)
	_, err = po.queries.insertOrderItem(ctx, itemsParams)
	if err != nil {
		return nil, err
	}
	order.Items = items
	return order, nil
}

func insertItemParamsFrom(orderID int64, items []models.OrderItem) []insertOrderItemParams {
	itemsParams := make([]insertOrderItemParams, 0, len(items))
	for _, it := range items {
		params := insertOrderItemParams{OrderID: orderID, SkuID: it.SkuId, Count: int32(it.Count)}
		itemsParams = append(itemsParams, params)
	}
	return itemsParams
}

// Save сохраняет заказ в БД в PostgreSQL. Изменение позиций заказа не предусмотрено
func (po *Orders) Save(ctx context.Context, order *models.Order) error {
	params := updateOrderParams{
		ID:               order.Id(),
		Status:           orderStatusToDTO(order.Status),
		AreItemsReserved: order.IsItemsReserved,
	}
	return po.queries.updateOrder(ctx, params)
}

func orderStatusToDTO(os models.OrderStatus) OrderStatus {
	switch os {
	case models.New:
		return OrderStatusNew
	case models.AwaitingPayment:
		return OrderStatusAwaitingPayment
	case models.Failed:
		return OrderStatusFailed
	case models.Payed:
		return OrderStatusPayed
	case models.Cancelled:
		return OrderStatusCancelled
	default:
		return OrderStatusUndefined
	}
}

func orderStatusFromDTO(s OrderStatus) models.OrderStatus {
	switch s {
	case OrderStatusNew:
		return models.New
	case OrderStatusAwaitingPayment:
		return models.AwaitingPayment
	case OrderStatusFailed:
		return models.Failed
	case OrderStatusPayed:
		return models.Payed
	case OrderStatusCancelled:
		return models.Cancelled
	default:
		return models.OrderStatus(0)
	}
}

// Load загружает информацию о заказе из БД в PostgreSQL, производя SELECT FOR UPDATE
func (po *Orders) Load(ctx context.Context, orderID int64) (*models.Order, error) {
	order, err := po.loadOrderRowWithoutItems(ctx, orderID)
	if err != nil {
		return nil, err
	}
	orderItems, err := po.readItemsForOrder(ctx, orderID)
	if err != nil {
		return nil, err
	}
	order.Items = orderItems
	return order, nil
}

func (po *Orders) loadOrderRowWithoutItems(ctx context.Context, orderID int64) (*models.Order, error) {
	row, err := po.queries.selectOrder(ctx, orderID)
	if err != nil {
		return nil, err
	}
	order := models.NewOrder(row.UserID, orderID)
	order.Status = orderStatusFromDTO(row.Status)
	order.IsItemsReserved = row.AreItemsReserved
	return order, nil
}

func (po *Orders) readItemsForOrder(ctx context.Context, orderID int64) ([]models.OrderItem, error) {
	rows, err := po.queries.selectOrderItems(ctx, orderID)
	if err != nil {
		return nil, err
	}
	items := make([]models.OrderItem, 0, len(rows))
	for _, it := range rows {
		item := models.OrderItem{SkuId: it.SkuID, Count: uint16(it.Count)}
		items = append(items, item)
	}
	return items, nil
}
