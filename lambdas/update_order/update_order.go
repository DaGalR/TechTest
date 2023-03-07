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
)

func UpdateOrderHandler(ctx context.Context, order *dto.UpdateOrderRequest) (events.APIGatewayProxyResponse, error){
	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err!=nil{
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error initializing DynamoDB client"} , err
	}
	client := dynamodb.NewFromConfig(config)
	_, err = adapters.UpdateOrderStatusDynamo(order.OrderID, order.NewStatus, client)
	if err != nil{
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: fmt.Sprintf("Could not update order with ID: %s because of: %s", order.OrderID, err.Error())}, err
	}
	return events.APIGatewayProxyResponse{StatusCode: 200, Body: fmt.Sprintf("Updated order status with ID: %s to ready for shipping",order.OrderID)}, nil
}

func main() {
	lambda.Start(UpdateOrderHandler)
}