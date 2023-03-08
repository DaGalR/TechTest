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

func UpdateOrderHandler(ctx context.Context, order *dto.UpdateOrderRequest) (events.APIGatewayProxyResponse, error){
	/*config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err!=nil{
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error initializing DynamoDB client"} , err
	}
	client := dynamodb.NewFromConfig(config)*/
	err := myAPI.UpdateOrderStatus(order.OrderID, order.NewStatus)
	if err != nil{
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: fmt.Sprintf("Could not update order with ID: %s because of: %s", order.OrderID, err.Error())}, err
	}
	return events.APIGatewayProxyResponse{StatusCode: 200, Body: fmt.Sprintf("Updated order status with ID: %s to ready for shipping",order.OrderID)}, nil
}

func main() {
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
	lambda.Start(UpdateOrderHandler)
}