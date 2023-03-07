package main

import (
	"context"
	"fmt"
	"os"

	"techtest/adapters"
	"techtest/dto"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func createOrderDynamo(order *dto.CreateOrderRequest, client *dynamodb.Client) error{

	err := adapters.CreateOrderDynamo(order, client)
	if err != nil{
		return fmt.Errorf("The order already exists in Dynamo")
	}
	return nil
}

func setOrderReadyForShipping(order *dto.CreateOrderRequest, client *dynamodb.Client) error{
	_, err := adapters.UpdateOrderStatusDynamo(order.OrderID, "Ready for shipping", client)
	if err != nil{
		return fmt.Errorf("The order could not be updated")
	}
	return nil
}

func OrdersHandler(ctx context.Context, order dto.CreateOrderRequest) (events.APIGatewayProxyResponse, error) {
	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err!=nil{
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error initializing DynamoDB client"} , err
	}
	sqsClient := sqs.NewFromConfig(config)
	client := dynamodb.NewFromConfig(config)
	err = createOrderDynamo(&order, client)
	if err != nil{
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: err.Error()}, err
	}
	var order_created_event dto.CreateOrderEvent
	order_created_event.OrderID = order.OrderID
	order_created_event.TotalPrice = order.TotalPrice
	adapters.SendOrderCreatedEvent("Order_Created", &order_created_event, sqsClient)
	if err != nil{
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: fmt.Sprintf("The order has been created but no event was sent to SQS: %s", err.Error())}, err
	}
	return events.APIGatewayProxyResponse{StatusCode: 200, Body: fmt.Sprintf("Order created with: %s USER: %s ITEM: %s QUANT: %d PRICE: %.2f", order.OrderID, order.UserID, order.Item, order.Quantity, order.TotalPrice)}, nil

}

func main(){
	lambda.Start(OrdersHandler)
}