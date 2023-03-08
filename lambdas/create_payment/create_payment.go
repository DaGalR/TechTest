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
func CreatePaymentHandler(ctx context.Context, payment dto.CreatePaymentRequest) (events.APIGatewayProxyResponse, error) {
	err:= myAPI.CreatePayment(&payment)
	if err != nil{
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: fmt.Sprintf("Error creating payment here's why: %s", err.Error())}, err
	}
	return events.APIGatewayProxyResponse{StatusCode: 200, Body: fmt.Sprintf("PAYMENT CREATED WITH ORDER: %s AND STATUS: %s", payment.OrderID, payment.Status)}, nil
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
	lambda.Start(CreatePaymentHandler)
}