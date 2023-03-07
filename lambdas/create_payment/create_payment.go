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

/*
func createPaymentDynamo(payment *dto.CreatePaymentRequest, client *dynamodb.Client) error{
	err := adapters.CreatePaymentDynamo(payment,client)
	if err != nil{
		return err
	}
	return nil
}
*/
var myAPI *entrypoints.API
func CreatePaymentHandler(ctx context.Context, payment dto.CreatePaymentRequest) (events.APIGatewayProxyResponse, error) {
	/*
	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err!=nil{
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error initializing DynamoDB client"} , err
	}
	sqsClient := sqs.NewFromConfig(config)
	client := dynamodb.NewFromConfig(config)
	fmt.Print("CREATING PAYMENT\n")
	err = createPaymentDynamo(&payment, client)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, err
	}
	adapters.SendPaymentCreatedEvent("Order_Completed", payment.OrderID, sqsClient)
	if err != nil{
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: fmt.Sprintf("The payment has been created but no event was sent to SQS: %s", err.Error())}, err
	}
	return events.APIGatewayProxyResponse{StatusCode: 200, Body: fmt.Sprintf("PAYMENT CREATED WITH ORDER: %s AND STATUS: %s", payment.OrderID, payment.Status)}, nil
	*/
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