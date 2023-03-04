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
	client := dynamodb.NewFromConfig(config)
	fmt.Print("CREATING PAYMENT\n")
	err = createPaymentDynamo(&payment, client)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, err
	}
	return events.APIGatewayProxyResponse{StatusCode: 200, Body: fmt.Sprintf("PAYMENT CREATED WITH ORDER: %s PRICE: %.2f", payment.OrderID, payment.TotalPrice)}, nil
}

func main(){
	lambda.Start(CreatePaymentHandler)
}