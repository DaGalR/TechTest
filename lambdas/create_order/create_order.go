package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"techtest/adapters"
	"techtest/dto"
)

func createOrderDynamo(order *dto.CreateOrderRequest, client *dynamodb.Client) error{

	err := adapters.CreateOrderDynamo(order, client)
	if err != nil{
		return fmt.Errorf("The order already exists in Dynamo")
	}
	return nil
}


func CreateOrderHandler(ctx context.Context, order dto.CreateOrderRequest) (events.APIGatewayProxyResponse, error) {
	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err!=nil{
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error initializing DynamoDB client"} , err
	}
	client := dynamodb.NewFromConfig(config)
	err = createOrderDynamo(&order, client)
	if err != nil{
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, err
	}
	return events.APIGatewayProxyResponse{StatusCode: 200, Body: fmt.Sprintf("Order created with: %s USER: %s ITEM: %s QUANT: %d PRICE: %.2f", order.OrderID, order.UserID, order.Item, order.Quantity, order.TotalPrice)}, nil
}
func main(){
	lambda.Start(CreateOrderHandler)
}