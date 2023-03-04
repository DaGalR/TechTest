package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)
const pkPayment = "PAYMENT"
const skPayment = "ID#%s"
type CreatePaymentRequest struct {
	OrderID    string `json:"order_id"`
	TotalPrice float32  `json:"total_price"`
}
func parseBodyToPaymentRequest(strObject string) (*CreatePaymentRequest, error){
	b := []byte(strObject)
	var responseData CreatePaymentRequest
	err := json.Unmarshal(b, &responseData)
	if err != nil{
		fmt.Printf("There was an error parsing the request body: %s", err.Error())
		return &CreatePaymentRequest{}, err
	}
	return &responseData,nil
}
func createPaymentDynamo(payment *CreatePaymentRequest, client *dynamodb.Client) error{
	_, err := client.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		ConditionExpression: aws.String("attribute_not_exists(SK)"),
		Item: map[string]types.AttributeValue{
			"PK":&types.AttributeValueMemberS{Value: pkPayment},
			"SK":&types.AttributeValueMemberS{Value: fmt.Sprintf(skPayment,payment.OrderID)},
			"totalPrice":&types.AttributeValueMemberN{Value: fmt.Sprint(payment.TotalPrice)},
		},
	})
	if err != nil && strings.Contains(err.Error(), "ConditionalCheckFailedException"){
		return fmt.Errorf("The payment already exists in Dynamo")
	}
	return nil
}
func CreatePaymentHandler(ctx context.Context, req CreatePaymentRequest) (events.APIGatewayProxyResponse, error) {
	/*payment,err := parseBodyToPaymentRequest(req.Body)
	if err != nil{
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: fmt.Sprintf("Error parsing request body: %s", err.Error())}, err
	}
	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err!=nil{
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error initializing DynamoDB client"} , err
	}
	client := dynamodb.NewFromConfig(config)
	err = createPaymentDynamo(payment, client)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, err
	}*/
	return events.APIGatewayProxyResponse{StatusCode: 200, Body: fmt.Sprintf("ORDER: %s PRICE: %.2f", req.OrderID, req.TotalPrice)}, nil
}

func main(){
	lambda.Start(CreatePaymentHandler)
}