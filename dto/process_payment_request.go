package dto

type ProcessPaymentRequest struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}