package converter

import (
	"github.com/samber/lo"
	"route256.ozon.ru/project/loms/internal/models"
	v1 "route256.ozon.ru/project/loms/internal/pkg/api/loms/v1"
)

// OrdersAllOrdersInfoResponse конвертирует []*models.Order в AllOrdersInfoResponse
func OrdersAllOrdersInfoResponse(orders []*models.Order) *v1.AllOrdersInfoResponse {
	dtoOrders := &v1.AllOrdersInfoResponse{}
	dtoOrders.OrdersInfo = lo.Map(orders, func(item *models.Order, _ int) *v1.OrderInfoWithID {
		orderInfoWithID := &v1.OrderInfoWithID{}
		orderInfoWithID.Id = IDToOrderID(item.Id())
		orderInfoWithID.OrdersInfo = OrderToOrderInfoResponse(item)
		return orderInfoWithID
	})
	return dtoOrders
}
