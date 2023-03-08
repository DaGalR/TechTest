package service

import (
	"fmt"
	"techtest/domain"
	"techtest/dto"
)
type HTTPPostUpdateOrder interface{
	CallUpdateOrdersService(string,string) error
}
type Domain interface {
	//Dynamo methods
	GetOrder(string)(*dto.CreateOrderRequest, error)
	CreateOrder(*dto.CreateOrderRequest) error
	UpdateOrderStatus(string, string,) error
	CreatePayment(*dto.CreatePaymentRequest) error
	//SQS Methods
	SendOrderCreatedEvent(string, *dto.CreateOrderEvent)(*string, error)
	SendPaymentCreatedEvent(string, string)(*string, error)
}

type Service struct {
	txDomain Domain
	repository domain.Repository
	httpPostUpdateOrder HTTPPostUpdateOrder
}

func NewService(txDomain Domain, repository domain.Repository, httpPostUpdateOrder HTTPPostUpdateOrder) *Service{
	return &Service{
		txDomain: txDomain,
		repository: repository,
		httpPostUpdateOrder: httpPostUpdateOrder,
	}
}

func (s *Service) CreateOrder(order *dto.CreateOrderRequest) error {
	err := s.txDomain.CreateOrder(order)
	if err != nil{
		return err
	}
	var msgID *string
	var order_created_event dto.CreateOrderEvent
	order_created_event.OrderID = order.OrderID
	order_created_event.TotalPrice = order.TotalPrice
	msgID, err = s.txDomain.SendOrderCreatedEvent("Order_Created",&order_created_event)
	if err != nil{
		return fmt.Errorf("The order with ID %s has been created but the order created event could not be sent: %s",order.OrderID,err.Error())
	}
	fmt.Printf("Order created and sent event with message ID: %s",*msgID)
	return nil
}

func (s *Service) CreatePayment(payment *dto.CreatePaymentRequest) error{
	err := s.txDomain.CreatePayment(payment)
	if err != nil{
		return err
	}
	var msgID *string
	msgID, err = s.txDomain.SendPaymentCreatedEvent("Order_Completed", payment.OrderID)
	if err != nil{
		return fmt.Errorf("The payment for order ID :%s was created but could not send the order complete event: %s", payment.OrderID, err.Error())
	}
	fmt.Printf("Payment created and sent event with message ID: %s",*msgID)
	return nil
}

func (s *Service) CallUpdateOrdersService(orderID, newStatus string) error{
	err := s.httpPostUpdateOrder.CallUpdateOrdersService(orderID, newStatus)
	if err != nil{
		return err
	}
	return nil
}

func (s *Service) UpdateOrderStatus(orderID string, newStatus string)error{
	err := s.txDomain.UpdateOrderStatus(orderID,newStatus)
	if err != nil{
		return err
	}
	return nil
}
