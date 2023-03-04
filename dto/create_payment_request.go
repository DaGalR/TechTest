package dto

type CreatePaymentRequest struct {
	OrderID    string `json:"order_id"`
	TotalPrice int64  `json:"total_price"`
}
