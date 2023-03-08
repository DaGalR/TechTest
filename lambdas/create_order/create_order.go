package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"techtest/adapters"
	"techtest/domain"
	"techtest/dto"
	"techtest/entrypoints"
	services "techtest/service"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

var myAPI *entrypoints.API
func OrdersHandler(ctx context.Context, order dto.CreateOrderRequest) (events.APIGatewayProxyResponse, error) {
	err:= myAPI.CreateOrder(&order)
	if err != nil{
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: fmt.Sprintf("Error creating order here's why: %s", err.Error())}, err
	}
	return events.APIGatewayProxyResponse{StatusCode: 200, Body: fmt.Sprintf("Order created with: %s USER: %s ITEM: %s QUANT: %d PRICE: %.2f", order.OrderID, order.UserID, order.Item, order.Quantity, order.TotalPrice)}, nil
}

func main(){
	
	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err!=nil{
		fmt.Errorf("Error creating config")
	}
	sqsClient := sqs.NewFromConfig(config)
	dynamoClient := dynamodb.NewFromConfig(config)
	httpClient := adapters.NewHTTPClient(&http.Client{})
	repository := struct{
		*adapters.DynamoAdapter
		*adapters.SQSAdapter
		
	}{
		adapters.NewDynamoAdapter(dynamoClient),
		adapters.NewSQSAdapter(sqsClient),
	}
	txDomain := domain.New(repository)
	service := services.NewService(txDomain, repository,httpClient)
	myAPI = entrypoints.NewAPI(service)
	lambda.Start(OrdersHandler)
}