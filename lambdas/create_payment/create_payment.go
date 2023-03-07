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

func createPaymentDynamo(payment *dto.CreatePaymentRequest, client *dynamodb.Client) error{
	err := adapters.CreatePaymentDynamo(payment,client)
	if err != nil{
		return err
	}
	return nil
}

func CreatePaymentHandler(ctx context.Context, payment dto.CreatePaymentRequest) (events.APIGatewayProxyResponse, error) {

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
}

func main(){
	lambda.Start(CreatePaymentHandler)
}