package dto

type UpdateOrderRequest struct {
	OrderID   string `json:"order_id"`
	NewStatus string `json:"new_status"`
}