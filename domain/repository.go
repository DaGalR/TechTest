package domain

import "techtest/dto"

type Repository interface {
	//Dynamo methods
	GetOrder(string) (*dto.CreateOrderRequest, error)
	CreateOrder(*dto.CreateOrderRequest) error
	UpdateOrderStatus(string, string) error
	CreatePayment(*dto.CreatePaymentRequest) error
	//SQS Methods
	SendOrderCreatedEvent(string, *dto.CreateOrderEvent)(*string, error)
	SendPaymentCreatedEvent(string, string)(*string, error)
}