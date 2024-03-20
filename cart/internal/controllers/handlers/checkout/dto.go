package checkout

type checkoutRequest struct {
	User *int64 `json:"user"`
}

type checkoutResponse struct {
	OrderId int64 `json:"orderID"`
}
