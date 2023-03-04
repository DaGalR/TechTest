package dto

type CreateOrderRequest struct {
	OrderID    string `json:"order_id"`
	UserID     string `json:"user_id"`
	Item       string `json:"item"`
	Quantity   int    `json:"quantity"`
	TotalPrice int64  `json:"total_price"`
}
