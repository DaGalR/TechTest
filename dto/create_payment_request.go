package dto

type CreatePaymentRequest struct {
	OrderID    string  `json:"order_id"`
	TotalPrice float32 `json:"total_price"`
}
