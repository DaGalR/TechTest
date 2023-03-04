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
	return nil
}

func startSQSLambda() {
	lambda.Start(sqsHandler)
}