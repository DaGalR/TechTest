package domain

import "techtest/dto"

type TxDomain struct {
	repository Repository
}

func New(repository Repository) *TxDomain {
	return &TxDomain{
		repository: repository,
	}
}

func (d *TxDomain) GetOrder(orderID string)  (*dto.CreateOrderRequest, error){
	res, err := d.repository.GetOrder(orderID)
	if err != nil{
		return &dto.CreateOrderRequest{}, err
	}
	return res,nil
}

func (d *TxDomain) CreateOrder(order *dto.CreateOrderRequest) error {
	err := d.repository.CreateOrder(order)
	if err != nil{
		return err
	}
	return nil
}
func (d *TxDomain) UpdateOrderStatus(orderID string, newStatus string) (map[string]map[string]interface{}, error){
	res, err := d.repository.UpdateOrderStatus(orderID,newStatus)
	if err != nil{
		return res, err
	}
	return res,nil
}
func (d *TxDomain) CreatePayment(payment *dto.CreatePaymentRequest) error{
	err := d.repository.CreatePayment(payment)
	if err != nil{
		return err
	}
	return nil
}

func (d *TxDomain) SendOrderCreatedEvent(body string, attributes *dto.CreateOrderEvent) (*string, error){
	var msgID *string
	msgID, err:=d.repository.SendOrderCreatedEvent(body, attributes)
	if err != nil{
		return nil, err 
	}
	return msgID,nil
}
func (d *TxDomain) SendPaymentCreatedEvent(body, orderID string) (*string, error){
	var msgID *string
	msgID, err:=d.repository.SendPaymentCreatedEvent(body, orderID)
	if err != nil{
		return nil, err 
	}
	return msgID,nil
}