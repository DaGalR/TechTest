package dto

type CreateOrderRequest struct {
	Operation  string  `json:"operation"`
	OrderID    string  `json:"order_id"`
	UserID     string  `json:"user_id"`
	Item       string  `json:"item"`
	Quantity   int     `json:"quantity"`
	TotalPrice float32 `json:"total_price"`
}
