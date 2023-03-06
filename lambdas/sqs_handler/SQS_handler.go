package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func sqsHandler(ctx context.Context, events events.SQSEvent) error{
	for _, message := range events.Records {
		fmt.Printf("The message %s for event source %s = %s \n", message.MessageId, message.EventSource, message.Body);
	}
	/*
	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err!=nil{
		fmt.Printf("Error creating client: %s", err.Error())
	}
	sqsClient := sqs.NewFromConfig(config)
	adapters.ReceiveMessage(sqsClient)*/
	return nil
}

func main() {
	lambda.Start(sqsHandler)
}