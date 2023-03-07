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

func SendOrderCreatedEvent(body string, attributes *dto.CreateOrderEvent, sqsClient *sqs.Client) (*string, error) {
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
	res, err := sqsClient.SendMessage(context.Background(), msgInput)
	if err!=nil{
		return nil, fmt.Errorf("There was an error sending the Order Created Event message: %s", err.Error())
	}
	return res.MessageId, nil
}

func SendPaymentCreatedEvent(body , orderID string, sqsClient *sqs.Client) (*string, error) {
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
	res, err := sqsClient.SendMessage(context.Background(), msgInput)
	if err!=nil{
		return nil, fmt.Errorf("There was an error sending the Order Completed Event message: %s", err.Error())
	}
	return res.MessageId, nil
}

func ReceiveMessage(sqsClient *sqs.Client) error{
	queueInput := &sqs.ReceiveMessageInput{
		MessageAttributeNames: []string{
			string(types.QueueAttributeNameAll),
		},
		QueueUrl:            aws.String(os.Getenv("SQS_QUEUE")),
		MaxNumberOfMessages: 1,
		VisibilityTimeout:   0,
	}
	msgResult, err := sqsClient.ReceiveMessage(context.TODO(), queueInput)
	if err != nil {
		return fmt.Errorf("Got an error receiving messages: %s", err.Error())
	}

	if msgResult.Messages != nil {
		fmt.Println("Message ID:     " + *msgResult.Messages[0].MessageId)
		fmt.Println("Message Handle: " + *msgResult.Messages[0].ReceiptHandle)
	} else {
		fmt.Println("No messages found")
	}
	return nil
}