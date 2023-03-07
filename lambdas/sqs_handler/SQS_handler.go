package main

import (
	"context"
	"fmt"
	"techtest/adapters"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func sqsHandler(ctx context.Context, sqsEvents events.SQSEvent) (events.APIGatewayProxyResponse, error){
	for _, message := range sqsEvents.Records {
		fmt.Printf("The message %s for event source %s = %s \n", message.MessageId, message.EventSource, message.Body);
		if message.Body == "Order_Completed"{
			orderid := *message.MessageAttributes["OrderID"].StringValue
			fmt.Printf("Recevied Order_Complete event with orderID: %s\n", orderid)
			err := adapters.CallUpdateOrdersService(orderid)
			if err != nil{
				return events.APIGatewayProxyResponse{StatusCode: 400, Body: fmt.Sprintf("Error sending request to API: %s", err.Error())}, err
			}
			return events.APIGatewayProxyResponse{StatusCode: 200, Body: fmt.Sprintf("Succesfuly sent request to API")}, nil
		}
	}
	return events.APIGatewayProxyResponse{StatusCode: 200, Body: fmt.Sprintf("Events processed")}, nil
}

func main() {
	lambda.Start(sqsHandler)
}