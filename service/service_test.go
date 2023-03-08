package service

import (
	"fmt"
	"techtest/dto"
	"techtest/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	mockDomain := new(mocks.Domain)
	mockRepo := new(mocks.Repository)
	mockHTTPClient := new(mocks.HTTPPostUpdateOrder)
	service:=NewService(mockDomain,mockRepo,mockHTTPClient)
	mockDomain.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	assert.NotNil(t,service)
}

func TestCreateOrder(t *testing.T){
	order := dto.CreateOrderRequest{
		OrderID: "01",
		Item: "Item",
		UserID: "TestUser",
		Quantity: 2,
		TotalPrice: 4.3,
	}
	order_created_event:=dto.CreateOrderEvent{
		OrderID: order.OrderID,
		TotalPrice: order.TotalPrice,
	}
	const body = "Order_Created"
	/*t.Run("Create order works", func(t *testing.T) {
		mockDomain := new(mocks.Domain)
		mockDomain.On("CreateOrder",&order).Return(nil)
		mockDomain.On("SendOrderCreatedEvent",body,&order_created_event).Return(nil,nil)
		mockRepo := new(mocks.Repository)
		mockHTTPClient := new(mocks.HTTPPostUpdateOrder)
		service:=NewService(mockDomain,mockRepo,mockHTTPClient)
		err:=service.CreateOrder(&order)
		mockDomain.AssertExpectations(t)
		assert.NoError(t,err)
	})*/
	t.Run("Create order fails", func(t *testing.T) {
		mockDomain := new(mocks.Domain)
		mockDomain.On("CreateOrder",&order).Return(nil)
		mockDomain.On("SendOrderCreatedEvent",body,&order_created_event).Return(nil,fmt.Errorf("There was an error sending the Order Created Event message: %s", "err"))
		mockRepo := new(mocks.Repository)
		mockHTTPClient := new(mocks.HTTPPostUpdateOrder)
		service:=NewService(mockDomain,mockRepo,mockHTTPClient)
		err:=service.CreateOrder(&order)
		mockDomain.AssertExpectations(t)
		assert.NotNil(t,err)
		assert.EqualError(t,err,"There was an error sending the Order Created Event message: fail")
	})
}