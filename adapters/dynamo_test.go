package adapters

import (
	"context"
	"fmt"
	"os"
	"techtest/dto"
	"techtest/mocks"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
)

func TestNewDynamoAdapter(t *testing.T){
	mockDynamoDBClient:=new(mocks.DynamoClient)
	adapter := NewDynamoAdapter(mockDynamoDBClient)
	mockDynamoDBClient.AssertExpectations(t)
	assert.NotNil(t, adapter)
}
func TestCreateOrderDynamo(t *testing.T){
	const pkOrder = "ORDER"
	const skOrder = "ID#%s"
	order := &dto.CreateOrderRequest{
		OrderID: "01",
		Item: "Item",
		UserID: "TestUser",
		Quantity: 2,
		TotalPrice: 4.3,
	} 
	t.Run("Add order correctly", func(t *testing.T){
		mockDynamoClient:=new(mocks.DynamoClient)
		mockDynamoClient.On("PutItem", context.TODO(), &dynamodb.PutItemInput{
			TableName: aws.String(os.Getenv("TABLE_NAME")),
			ConditionExpression:aws.String("attribute_not_exists(SK)"),
			Item: map[string]types.AttributeValue{
			"PK":         &types.AttributeValueMemberS{Value: pkOrder},
			"SK":         &types.AttributeValueMemberS{Value: fmt.Sprintf(skOrder, order.OrderID)},
			"userID":     &types.AttributeValueMemberS{Value: order.UserID},
			"item":       &types.AttributeValueMemberS{Value: order.Item},
			"quantity":   &types.AttributeValueMemberN{Value: fmt.Sprint(order.Quantity)},
			"totalPrice": &types.AttributeValueMemberN{Value: fmt.Sprint(order.TotalPrice)},
			"status":     &types.AttributeValueMemberS{Value: "Incomplete"},
		},
		}).Return(&dynamodb.PutItemOutput{},nil)
		adapter := NewDynamoAdapter(mockDynamoClient)
		err:=adapter.CreateOrder(order)
		mockDynamoClient.AssertExpectations(t)
		assert.NoError(t,err)
	})
	t.Run("Order exists in DB, so doesnt create it", func(t *testing.T){
		mockDynamoClient:=new(mocks.DynamoClient)
		mockDynamoClient.On("PutItem", context.TODO(), &dynamodb.PutItemInput{
			TableName: aws.String(os.Getenv("TABLE_NAME")),
			ConditionExpression:aws.String("attribute_not_exists(SK)"),
			Item: map[string]types.AttributeValue{
			"PK":         &types.AttributeValueMemberS{Value: pkOrder},
			"SK":         &types.AttributeValueMemberS{Value: fmt.Sprintf(skOrder, order.OrderID)},
			"userID":     &types.AttributeValueMemberS{Value: order.UserID},
			"item":       &types.AttributeValueMemberS{Value: order.Item},
			"quantity":   &types.AttributeValueMemberN{Value: fmt.Sprint(order.Quantity)},
			"totalPrice": &types.AttributeValueMemberN{Value: fmt.Sprint(order.TotalPrice)},
			"status":     &types.AttributeValueMemberS{Value: "Incomplete"},
		},
		}).Return(&dynamodb.PutItemOutput{},fmt.Errorf("ConditionalCheckFailedException"))
		adapter := NewDynamoAdapter(mockDynamoClient)
		err:=adapter.CreateOrder(order)
		mockDynamoClient.AssertExpectations(t)
		assert.EqualError(t,err,"The order already exists in Dynamo")
	})
}

func TestCreatePaymentDynamo(t *testing.T){
	const pkPayment = "PAYMENT"
	const skPayment = "ID#%s"
	const pkOrder = "ORDER"
	const skOrder = "ID#%s"

	payment := &dto.CreatePaymentRequest{
		OrderID: "01",
		Status: "Complete",
	} 
	order := &dto.CreateOrderRequest{
		OrderID: "01",
		Item: "Item",
		UserID: "TestUser",
		Quantity: 2,
		TotalPrice: 4.3,
	} 
	t.Run("Payment wont be added because there was an error retrieving the order", func(t *testing.T){
		mockDynamoClient:=new(mocks.DynamoClient)
		mockDynamoClient.On("GetItem", context.TODO(), &dynamodb.GetItemInput{
			TableName: aws.String(os.Getenv("TABLE_NAME")),
			Key: map[string]types.AttributeValue{
				"PK": &types.AttributeValueMemberS{Value: pkOrder},
				"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf(skOrder, order.OrderID)},
			},
		}).Return(&dynamodb.GetItemOutput{
			Item: nil,
		}, fmt.Errorf("Could not verify if order with id: " + payment.OrderID + " exists because of internal server error"))
		
		adapter := NewDynamoAdapter(mockDynamoClient)
		err:=adapter.CreatePayment(payment)
		mockDynamoClient.AssertExpectations(t)
		assert.EqualError(t,err, fmt.Sprintf("Could not verify if order with id: " + payment.OrderID + " exists because of internal server error"))
	})
	/*t.Run("Order does not exist in DB, so payment is not created", func(t *testing.T){
		mockDynamoClient:=new(mocks.DynamoClient)
		mockDynamoClient.On("GetItem", context.TODO(), &dynamodb.GetItemInput{
			TableName: aws.String(os.Getenv("TABLE_NAME")),
			Key: map[string]types.AttributeValue{
				"PK": &types.AttributeValueMemberS{Value: pkOrder},
				"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf(skOrder, order.OrderID)},
			},
		}).Return(&dynamodb.GetItemOutput{
			Item:nil,
		}, nil)
		adapter := NewDynamoAdapter(mockDynamoClient)
		err:=adapter.CreatePayment(payment)
		mockDynamoClient.AssertExpectations(t)
		assert.EqualError(t,err, fmt.Sprintf("There is no order with id: "+ payment.OrderID))
	})*/
	t.Run("Payment already exists in dynamo", func(t *testing.T) {
		mockDynamoClient:=new(mocks.DynamoClient)
		mockDynamoClient.On("GetItem", context.TODO(), &dynamodb.GetItemInput{
			TableName: aws.String(os.Getenv("TABLE_NAME")),
			Key: map[string]types.AttributeValue{
				"PK": &types.AttributeValueMemberS{Value: pkOrder},
				"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf(skOrder, order.OrderID)},
			},
		}).Return(&dynamodb.GetItemOutput{
			Item: map[string]types.AttributeValue{
				"PK": &types.AttributeValueMemberS{Value: pkOrder},
				"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf(skOrder, order.OrderID)},
				"OrderID": &types.AttributeValueMemberS{Value:"01"},
				"UserID": &types.AttributeValueMemberS{Value:"TestUser"},
				"Quantity": &types.AttributeValueMemberN{Value: "2"},
				"TotalPrice": &types.AttributeValueMemberN{Value:"4.3"},
			},
		},nil)
		mockDynamoClient.On("PutItem", context.TODO(), &dynamodb.PutItemInput{
			TableName: aws.String(os.Getenv("TABLE_NAME")),
			ConditionExpression: aws.String("attribute_not_exists(SK)"),
			Item: map[string]types.AttributeValue{
				"PK":&types.AttributeValueMemberS{Value: pkPayment},
				"SK":&types.AttributeValueMemberS{Value: fmt.Sprintf(skPayment,order.OrderID)},
				"status":&types.AttributeValueMemberS{Value: payment.Status},
			},
		}).Return(nil, fmt.Errorf("ConditionalCheckFailedException"))
		adapter := NewDynamoAdapter(mockDynamoClient)
		err:=adapter.CreatePayment(payment)
		mockDynamoClient.AssertExpectations(t)
		assert.EqualError(t,err,"The payment already exists in Dynamo")
	})

	t.Run("Payment was created succesfuly", func(t *testing.T) {
		mockDynamoClient:=new(mocks.DynamoClient)
		mockDynamoClient.On("GetItem", context.TODO(), &dynamodb.GetItemInput{
			TableName: aws.String(os.Getenv("TABLE_NAME")),
			Key: map[string]types.AttributeValue{
				"PK": &types.AttributeValueMemberS{Value: pkOrder},
				"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf(skOrder, order.OrderID)},
			},
		}).Return(&dynamodb.GetItemOutput{
			Item: map[string]types.AttributeValue{
				"PK": &types.AttributeValueMemberS{Value: pkOrder},
				"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf(skOrder, order.OrderID)},
				"OrderID": &types.AttributeValueMemberS{Value:"01"},
				"UserID": &types.AttributeValueMemberS{Value:"TestUser"},
				"Quantity": &types.AttributeValueMemberN{Value: "2"},
				"TotalPrice": &types.AttributeValueMemberN{Value:"4.3"},
			},
		},nil)
		mockDynamoClient.On("PutItem", context.TODO(), &dynamodb.PutItemInput{
			TableName: aws.String(os.Getenv("TABLE_NAME")),
			ConditionExpression: aws.String("attribute_not_exists(SK)"),
			Item: map[string]types.AttributeValue{
				"PK":&types.AttributeValueMemberS{Value: pkPayment},
				"SK":&types.AttributeValueMemberS{Value: fmt.Sprintf(skPayment,order.OrderID)},
				"status":&types.AttributeValueMemberS{Value: payment.Status},
			},
		}).Return(&dynamodb.PutItemOutput{}, nil)
		adapter := NewDynamoAdapter(mockDynamoClient)
		err:=adapter.CreatePayment(payment)
		mockDynamoClient.AssertExpectations(t)
		assert.NoError(t,err)
	})
}

func TestGetOrder(t *testing.T){
	const pkOrder = "ORDER"
	const skOrder = "ID#%s"
	order := &dto.CreateOrderRequest{
		OrderID: "01",
		Item: "Item",
		UserID: "TestUser",
		Quantity: 2,
		TotalPrice: 4.3,
	} 
	t.Run("Get order runs correctly", func(t *testing.T) {
		mockDynamoClient:= new(mocks.DynamoClient)
		mockDynamoClient.On("GetItem", context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pkOrder},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf(skOrder, order.OrderID)},
		},
		}).Return(&dynamodb.GetItemOutput{
			Item: map[string]types.AttributeValue{
				"PK": &types.AttributeValueMemberS{Value: pkOrder},
				"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf(skOrder, order.OrderID)},
				"OrderID": &types.AttributeValueMemberS{Value:"01"},
				"Item": &types.AttributeValueMemberS{Value: "Item"},
				"UserID": &types.AttributeValueMemberS{Value:"TestUser"},
				"Quantity": &types.AttributeValueMemberN{Value: "2"},
				"TotalPrice": &types.AttributeValueMemberN{Value:"4.3"},
			},
		},nil)
		adapter:= NewDynamoAdapter(mockDynamoClient)
		got,err:=adapter.GetOrder(order.OrderID)
		mockDynamoClient.AssertExpectations(t)
		assert.Equal(t, &dto.CreateOrderRequest{
			OrderID: "01",
			Item: "Item",
			UserID: "TestUser",
			Quantity: 2,
			TotalPrice: 4.3,
		},got)
		assert.NoError(t, err)
	})

	t.Run("Get order fails", func(t *testing.T) {
		mockDynamoClient:= new(mocks.DynamoClient)
		mockDynamoClient.On("GetItem", context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pkOrder},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf(skOrder, order.OrderID)},
		},
		}).Return(&dynamodb.GetItemOutput{
			Item: map[string]types.AttributeValue{
				"PK": &types.AttributeValueMemberS{Value: pkOrder},
				"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf(skOrder, order.OrderID)},
				"OrderID": &types.AttributeValueMemberS{Value:"01"},
				"UserID": &types.AttributeValueMemberS{Value:"TestUser"},
				"Quantity": &types.AttributeValueMemberN{Value: "2"},
				"TotalPrice": &types.AttributeValueMemberN{Value:"4.3"},
			},
		},fmt.Errorf("There was an error retrieving data from Dynamo"))
		adapter:= NewDynamoAdapter(mockDynamoClient)
		got,err:=adapter.GetOrder(order.OrderID)
		mockDynamoClient.AssertExpectations(t)
		assert.Equal(t, &dto.CreateOrderRequest{},got)
		assert.EqualError(t, err,"There was an error retrieving data from Dynamo" )
	})

		t.Run("Get order fails unmarshal", func(t *testing.T) {
		mockDynamoClient:= new(mocks.DynamoClient)
		mockDynamoClient.On("GetItem", context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pkOrder},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf(skOrder, order.OrderID)},
		},
		}).Return(&dynamodb.GetItemOutput{
			Item: map[string]types.AttributeValue{
				"PK": &types.AttributeValueMemberS{Value: pkOrder},
				"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf(skOrder, order.OrderID)},
				"OrderID": &types.AttributeValueMemberS{Value:"01"},
				"UserID": &types.AttributeValueMemberS{Value:"TestUser"},
				"Quantity": &types.AttributeValueMemberN{Value: "0"},
				"TotalPrice": &types.AttributeValueMemberS{Value:"4.3"},
			},
		},nil)
		adapter:= NewDynamoAdapter(mockDynamoClient)
		got,err:=adapter.GetOrder(order.OrderID)
		mockDynamoClient.AssertExpectations(t)
		assert.Equal(t, &dto.CreateOrderRequest{},got)
		assert.EqualError(t, err,"UnmarshalMap: unmarshal failed, cannot unmarshal string into Go value type float32")
	})
}

func TestUpdateOrderStatus(t *testing.T){
	const newStatus string ="Completed"
	const orderID string ="01"
	const pkOrder = "ORDER"
	const skOrder = "ID#%s"
	t.Run("Update fails server error", func(t *testing.T) {
		mockDynamoClient:= new(mocks.DynamoClient)
		update := expression.Set(expression.Name("status"), expression.Value(newStatus))
		expr, err := expression.NewBuilder().WithUpdate(update).Build()
		mockDynamoClient.On("UpdateItem", context.TODO(), &dynamodb.UpdateItemInput{
			TableName: aws.String(os.Getenv("TABLE_NAME")),
			Key: map[string]types.AttributeValue{
				"PK": &types.AttributeValueMemberS{Value: pkOrder},
				"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf(skOrder, orderID)},
			},
			ExpressionAttributeNames: expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			UpdateExpression: expr.Update(),
		}).Return(&dynamodb.UpdateItemOutput{},fmt.Errorf("Couldn't update order with id: %s", orderID))
		adapter := NewDynamoAdapter(mockDynamoClient)
		err=adapter.UpdateOrderStatus(orderID,newStatus)
		mockDynamoClient.AssertExpectations(t)
		assert.EqualError(t, err, fmt.Sprintf("Couldn't update order with id: %s", orderID))
	}) 
	t.Run("Update works", func(t *testing.T) {
		mockDynamoClient:= new(mocks.DynamoClient)
		update := expression.Set(expression.Name("status"), expression.Value(newStatus))
		expr, err := expression.NewBuilder().WithUpdate(update).Build()
		mockDynamoClient.On("UpdateItem", context.TODO(), &dynamodb.UpdateItemInput{
			TableName: aws.String(os.Getenv("TABLE_NAME")),
			Key: map[string]types.AttributeValue{
				"PK": &types.AttributeValueMemberS{Value: pkOrder},
				"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf(skOrder, orderID)},
			},
			ExpressionAttributeNames: expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			UpdateExpression: expr.Update(),
		}).Return(&dynamodb.UpdateItemOutput{
			Attributes: map[string]types.AttributeValue{
				"PK": &types.AttributeValueMemberS{Value: pkOrder},
				"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf(skOrder, orderID)},
				"OrderID": &types.AttributeValueMemberS{Value:"01"},
				"UserID": &types.AttributeValueMemberS{Value:"TestUser"},
				"Quantity": &types.AttributeValueMemberN{Value: "2"},
				"TotalPrice": &types.AttributeValueMemberN{Value:"4.3"},
				"Item": &types.AttributeValueMemberS{Value:"MockItem"},
				"Status": &types.AttributeValueMemberS{Value: newStatus},
			},
		},nil)
		adapter := NewDynamoAdapter(mockDynamoClient)
		err=adapter.UpdateOrderStatus(orderID,newStatus)
		mockDynamoClient.AssertExpectations(t)
		assert.NoError(t,err)
	})
}