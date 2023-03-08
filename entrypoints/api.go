package entrypoints

import "techtest/dto"

type API struct {
	service Service
}
type Service interface {
	CreateOrder(*dto.CreateOrderRequest) error
	UpdateOrderStatus(string, string) error
	CreatePayment(*dto.CreatePaymentRequest) error
	CallUpdateOrdersService(string,string) error
}

func NewAPI(service Service) *API{
	return &API{
		service: service,
	}
}

func (a *API) CreateOrder(order *dto.CreateOrderRequest) error{
	err := a.service.CreateOrder(order)
	if err != nil{
		return err
	}
	return nil
}

func (a *API) UpdateOrderStatus(orderID, newStatus string) error{
	err := a.service.UpdateOrderStatus(orderID, newStatus)
	if err != nil{
		return err
	}
	return nil
}

func (a *API) CreatePayment(payment *dto.CreatePaymentRequest) error{
	err := a.service.CreatePayment(payment)
	if err != nil{
		return err
	}
	return nil
}

func (a *API) CallUpdateOrdersService(orderID,newStatus string) error{
	err := a.service.CallUpdateOrdersService(orderID,newStatus)
	if err != nil{
		return err
	}
	return nil
}