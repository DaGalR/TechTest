package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"techtest/adapters"
	"techtest/domain"
	"techtest/entrypoints"
	services "techtest/service"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)
var myAPI *entrypoints.API
func sqsHandler(ctx context.Context, sqsEvents events.SQSEvent) (events.APIGatewayProxyResponse, error){
	for _, message := range sqsEvents.Records {
		fmt.Printf("The message %s for event source %s = %s \n", message.MessageId, message.EventSource, message.Body);
		//Event when a payment is done for an order, the order status is set to ready for shipping
		if message.Body == "Order_Completed"{
			orderid := *message.MessageAttributes["OrderID"].StringValue
			fmt.Printf("Recevied Order_Complete event with orderID: %s\n", orderid)
			err := myAPI.CallUpdateOrdersService(orderid, "Ready for shipping")
			if err != nil{
				return events.APIGatewayProxyResponse{StatusCode: 400, Body: fmt.Sprintf("Error sending request to API: %s", err.Error())}, err
			}
			return events.APIGatewayProxyResponse{StatusCode: 200, Body: fmt.Sprintf("Succesfuly sent request to API")}, nil
		}
	}
	return events.APIGatewayProxyResponse{StatusCode: 200, Body: fmt.Sprintf("Events processed")}, nil
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
	lambda.Start(sqsHandler)
}