package domain

import (
	"fmt"
	"techtest/dto"
	"techtest/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDomain(t *testing.T) {
	mockRepository := new(mocks.Repository)
	domain:=New(mockRepository)
	mockRepository.AssertExpectations(t)
	assert.NotNil(t, domain)
}

func TestGetOrder(t *testing.T){
	order := &dto.CreateOrderRequest{
		OrderID: "01",
		Item: "Item",
		UserID: "TestUser",
		Quantity: 2,
		TotalPrice: 4.3,
	} 
	t.Run("Get order works fine", func (t *testing.T)  {
		mockRepository:=new(mocks.Repository)
		mockRepository.On("GetOrder", order.OrderID).Return(&dto.CreateOrderRequest{
			OrderID: "01",
			Item: "Item",
			UserID: "TestUser",
			Quantity: 2,
			TotalPrice: 4.3,
		}, nil)
		domain:= New(mockRepository)
		got,err:=domain.GetOrder(order.OrderID)
		mockRepository.AssertExpectations(t)
		assert.NoError(t,err)
		assert.Equal(t,order,got)
	})
	t.Run("Get order fails", func (t *testing.T)  {
		mockRepository:=new(mocks.Repository)
		mockRepository.On("GetOrder", order.OrderID).Return(&dto.CreateOrderRequest{}, fmt.Errorf("There was an error retrieving data from Dynamo"))
		domain:= New(mockRepository)
		got,err:=domain.GetOrder(order.OrderID)
		mockRepository.AssertExpectations(t)
		assert.EqualError(t,err,"There was an error retrieving data from Dynamo")
		assert.Equal(t,&dto.CreateOrderRequest{},got)
	})
}

func TestCreateOrder(t *testing.T){
	order := dto.CreateOrderRequest{
		OrderID: "01",
		Item: "Item",
		UserID: "TestUser",
		Quantity: 2,
		TotalPrice: 4.3,
	} 
	t.Run("Create order works", func (t *testing.T)  {
		mockRepository:=new(mocks.Repository)
		mockRepository.On("CreateOrder", &order).Return(nil)
		domain:= New(mockRepository)
		err:=domain.CreateOrder(&order)
		mockRepository.AssertExpectations(t)
		assert.NoError(t,err)
	})

	t.Run("Create order fails", func (t *testing.T)  {
		mockRepository:=new(mocks.Repository)
		mockRepository.On("CreateOrder", &order).Return(fmt.Errorf("The order already exists in Dynamo"))
		domain:= New(mockRepository)
		err:=domain.CreateOrder(&order)
		mockRepository.AssertExpectations(t)
		assert.NotNil(t,err)
	})
}

func TestUpdateOrderStatus(t *testing.T){
	const orderID = "01"
	const newStatus = "Completed"

	t.Run("Update order works", func (t *testing.T) {
		mockRepository:=new(mocks.Repository)
		mockRepository.On("UpdateOrderStatus", orderID,newStatus).Return(nil)
		domain:= New(mockRepository)
		err:=domain.UpdateOrderStatus(orderID, newStatus)
		mockRepository.AssertExpectations(t)
		assert.NoError(t,err)
	})

	t.Run("Update order fails", func (t *testing.T) {
		mockRepository:=new(mocks.Repository)
		mockRepository.On("UpdateOrderStatus", orderID,newStatus).Return(fmt.Errorf("Couldn't update order with id: %s", orderID))
		domain:= New(mockRepository)
		err:=domain.UpdateOrderStatus(orderID, newStatus)
		mockRepository.AssertExpectations(t)
		assert.NotNil(t,err)
	})
}

func TestCreatePayment(t *testing.T){
	payment := dto.CreatePaymentRequest{
		OrderID: "01",
		Status: "Complete",
	}
	t.Run("Create payment fails", func (t *testing.T) {
		mockRepository:=new(mocks.Repository)
		mockRepository.On("CreatePayment", &payment).Return(fmt.Errorf("The payment already exists in Dynamo"))
		domain:= New(mockRepository)
		err:=domain.CreatePayment(&payment)
		mockRepository.AssertExpectations(t)
		assert.NotNil(t,err)
	})

	t.Run("Create payment works", func (t *testing.T) {
		mockRepository:=new(mocks.Repository)
		mockRepository.On("CreatePayment", &payment).Return(nil)
		domain:= New(mockRepository)
		err:=domain.CreatePayment(&payment)
		mockRepository.AssertExpectations(t)
		assert.NoError(t,err)
	})
}

func TestSendOrderCreatedEvent(t *testing.T){
	const body = "Order_Created"
	createOrderEvent := dto.CreateOrderEvent{
		OrderID: "01",
		TotalPrice: 3.2,
	}

	t.Run("Send order create event works", func(t *testing.T) {
		var res *string
		mockRepository:=new(mocks.Repository)
		mockRepository.On("SendOrderCreatedEvent", body, &createOrderEvent).Return(res,nil)
		domain:= New(mockRepository)
		res, err:=domain.SendOrderCreatedEvent(body,&createOrderEvent)
		mockRepository.AssertExpectations(t)
		assert.NoError(t,err)
	})

	t.Run("Send order create event fails", func(t *testing.T) {
		mockRepository:=new(mocks.Repository)
		mockRepository.On("SendOrderCreatedEvent", body, &createOrderEvent).Return(nil,fmt.Errorf("There was an error sending the Order Created Event message: %s", "dx"))
		domain:= New(mockRepository)
		_, err:=domain.SendOrderCreatedEvent(body,&createOrderEvent)
		mockRepository.AssertExpectations(t)
		assert.NotNil(t,err)
	})
}

func TestSendPaymentCreatedEvent(t *testing.T){
	const body = "Order_Complete"	
	const orderID =  "01"

	t.Run("Send payment created event works", func(t *testing.T) {
		var res *string
		mockRepository:=new(mocks.Repository)
		mockRepository.On("SendPaymentCreatedEvent", body, orderID).Return(res,nil)
		domain:= New(mockRepository)
		res, err:=domain.SendPaymentCreatedEvent(body,orderID)
		mockRepository.AssertExpectations(t)
		assert.NoError(t,err)
	})

	t.Run("Send payment created event fails", func(t *testing.T) {
		mockRepository:=new(mocks.Repository)
		mockRepository.On("SendPaymentCreatedEvent", body, orderID).Return(nil,fmt.Errorf("There was an error sending the Order Created Event message: %s", "dx"))
		domain:= New(mockRepository)
		_, err:=domain.SendPaymentCreatedEvent(body,orderID)
		mockRepository.AssertExpectations(t)
		assert.NotNil(t,err)
	})
}