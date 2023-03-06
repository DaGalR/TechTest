package dto

type CreatePaymentRequest struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}
