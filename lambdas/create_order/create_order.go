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
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)
const pkOrder = "ORDER"
const skOrder = "ID#%s"
type CreateOrderRequest struct {
	OrderID    string `json:"order_id"`
	UserID     string `json:"user_id"`
	Item       string `json:"item"`
	Quantity   int    `json:"quantity"`
	TotalPrice float32  `json:"total_price"`
}

func parseBodyToOrderObject(strObject string) (*CreateOrderRequest, error){
	b := []byte(strObject)
	var responseData CreateOrderRequest
	err := json.Unmarshal(b, &responseData)
	if err != nil{
		fmt.Printf("There was an error parsing the request body: %s", err.Error())
		return &CreateOrderRequest{}, err
	}
	return &responseData,nil
}

func createOrderDynamo(order *CreateOrderRequest, client *dynamodb.Client) error{
	_, err := client.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		ConditionExpression: aws.String("attribute_not_exists(SK)"),
		Item: map[string]types.AttributeValue{
			"PK":&types.AttributeValueMemberS{Value: pkOrder},
			"SK":&types.AttributeValueMemberS{Value: fmt.Sprintf(skOrder, order.OrderID)},
			"userID":&types.AttributeValueMemberS{Value: order.UserID},
			"item":&types.AttributeValueMemberS{Value: order.Item},
			"quantity":&types.AttributeValueMemberN{Value: fmt.Sprint(order.Quantity)},
			"totalPrice":&types.AttributeValueMemberN{Value: fmt.Sprint(order.TotalPrice)},
			"status": &types.AttributeValueMemberS{Value: "Incomplete"},
		},
	})
	if err != nil && strings.Contains(err.Error(), "ConditionalCheckFailedException"){
		return fmt.Errorf("The order already exists in Dynamo")
	}
	return nil
}


func CreateOrderHandler(ctx context.Context, order CreateOrderRequest) (events.APIGatewayProxyResponse, error) {
	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err!=nil{
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error initializing DynamoDB client"} , err
	}
	client := dynamodb.NewFromConfig(config)
	err = createOrderDynamo(&order, client)
	if err != nil{
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, err
	}
	return events.APIGatewayProxyResponse{StatusCode: 200, Body: fmt.Sprintf("Order created with: %s USER: %s ITEM: %s QUANT: %d PRICE: %.2f", order.OrderID, order.UserID, order.Item, order.Quantity, order.TotalPrice)}, nil
}
func main(){
	lambda.Start(CreateOrderHandler)
}