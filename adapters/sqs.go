package adapters

import (
	"context"
	"fmt"
	"os"
	"techtest/dto"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type SQSAdapter struct{
	sqsClient SQSClient
}
type SQSClient interface{
	SendMessage(context.Context, *sqs.SendMessageInput, ...func(*sqs.Options))(*sqs.SendMessageOutput, error)
}

func NewSQSAdapter(client SQSClient) *SQSAdapter{
	return &SQSAdapter{
		sqsClient: client,
	}
}

func (s *SQSAdapter) SendOrderCreatedEvent(body string, attributes *dto.CreateOrderEvent) (*string, error) {
	msgInput := &sqs.SendMessageInput{
		MessageAttributes: map[string]types.MessageAttributeValue{
			"OrderID": {
				DataType:    aws.String("String"),
				StringValue: aws.String(attributes.OrderID),
			},
			"Total_price": {
				DataType:    aws.String("Number"),
				StringValue: aws.String(fmt.Sprintf("%.2f",attributes.TotalPrice)),
			},
		},
		MessageBody: aws.String(body),
		QueueUrl:    aws.String(os.Getenv("SQS_QUEUE")),
	}
	res, err := s.sqsClient.SendMessage(context.Background(), msgInput)
	if err!=nil{
		return nil, fmt.Errorf("There was an error sending the Order Created Event message: %s", err.Error())
	}
	return res.MessageId, nil
}

func (s *SQSAdapter) SendPaymentCreatedEvent(body, orderID string) (*string, error) {
	msgInput := &sqs.SendMessageInput{
		MessageAttributes: map[string]types.MessageAttributeValue{
			"OrderID": {
				DataType:    aws.String("String"),
				StringValue: aws.String(orderID),
			},
		},
		MessageBody: aws.String(body),
		QueueUrl:    aws.String(os.Getenv("SQS_QUEUE")),
	}
	res, err := s.sqsClient.SendMessage(context.Background(), msgInput)
	if err!=nil{
		return nil, fmt.Errorf("There was an error sending the Order Completed Event message: %s", err.Error())
	}
	return res.MessageId, nil
}