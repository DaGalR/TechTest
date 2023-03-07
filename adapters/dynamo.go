package adapters

import (
	"context"
	"fmt"
	"os"
	"strings"
	"techtest/dto"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)
const pkPayment = "PAYMENT"
const skPayment = "ID#%s"
const pkOrder = "ORDER"
const skOrder = "ID#%s"

type DynamoAdapter struct{
	client DynamoClient
}
type DynamoClient interface{
	GetItem(context.Context, *dynamodb.GetItemInput, ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	PutItem(context.Context, *dynamodb.PutItemInput, ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	UpdateItem(context.Context, *dynamodb.UpdateItemInput, ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error)
}
func NewDynamoAdapter(client DynamoClient) *DynamoAdapter{
	return &DynamoAdapter{
		client: client,
	}
}
func (d *DynamoAdapter) CreateOrder(order *dto.CreateOrderRequest) error {
	_, err := d.client.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName:           aws.String(os.Getenv("TABLE_NAME")),
		ConditionExpression: aws.String("attribute_not_exists(SK)"),
		Item: map[string]types.AttributeValue{
			"PK":         &types.AttributeValueMemberS{Value: pkOrder},
			"SK":         &types.AttributeValueMemberS{Value: fmt.Sprintf(skOrder, order.OrderID)},
			"userID":     &types.AttributeValueMemberS{Value: order.UserID},
			"item":       &types.AttributeValueMemberS{Value: order.Item},
			"quantity":   &types.AttributeValueMemberN{Value: fmt.Sprint(order.Quantity)},
			"totalPrice": &types.AttributeValueMemberN{Value: fmt.Sprint(order.TotalPrice)},
			"status":     &types.AttributeValueMemberS{Value: "Incomplete"},
		},
	})
	if err != nil && strings.Contains(err.Error(), "ConditionalCheckFailedException") {
		return fmt.Errorf("The order already exists in Dynamo")
	}
	return nil
}

func (d *DynamoAdapter) CreatePayment(payment *dto.CreatePaymentRequest) error{
	order, err := d.GetOrder(payment.OrderID)
	if err != nil{
		return fmt.Errorf("Could not verify if order with id: " + payment.OrderID + " exists because of: " + err.Error())
	}
	if order == (&dto.CreateOrderRequest{}){
		return fmt.Errorf("There is no order with id: "+ payment.OrderID)
	}else{
		_, err = d.client.PutItem(context.Background(), &dynamodb.PutItemInput{
			TableName: aws.String(os.Getenv("TABLE_NAME")),
			ConditionExpression: aws.String("attribute_not_exists(SK)"),
			Item: map[string]types.AttributeValue{
				"PK":&types.AttributeValueMemberS{Value: pkPayment},
				"SK":&types.AttributeValueMemberS{Value: fmt.Sprintf(skPayment,payment.OrderID)},
				"status":&types.AttributeValueMemberS{Value: payment.Status},
			},
		})
		if err != nil && strings.Contains(err.Error(), "ConditionalCheckFailedException"){
			return fmt.Errorf("The payment already exists in Dynamo")
		}
	}
	return nil
}

func (d *DynamoAdapter) GetOrder(orderID string) (*dto.CreateOrderRequest, error) {
	order := dto.CreateOrderRequest{}
	fmt.Printf("LOOKING FOR ORDER WITH ID %s\n",orderID)
	data, err := d.client.GetItem(context.Background(), &dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pkOrder},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf(skOrder, orderID)},
		},
	},)
	fmt.Printf("commerce data %v\n", data.Item)
	if err!=nil{
		return &order, fmt.Errorf("There was an error retrieving data from Dynamo")
	}
	if data.Item == nil{
		return &order, fmt.Errorf("Order not found")
	}
	err = attributevalue.UnmarshalMap(data.Item, &order)
	if err != nil{
		return &order, fmt.Errorf("UnmarshalMap: %v", err)
	}
	return &order,nil
}

func (d *DynamoAdapter) UpdateOrderStatus(orderID string, newStatus string)(map[string]map[string]interface{}, error){
	update := expression.Set(expression.Name("status"), expression.Value(newStatus))
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	var response *dynamodb.UpdateItemOutput
	var attributeMap map[string]map[string]interface{}
	if err != nil{
		return attributeMap, fmt.Errorf("Couldn't build expression for order update: %v\n",err)
	}else{
		response, err = d.client.UpdateItem(context.Background(), &dynamodb.UpdateItemInput{
			TableName: aws.String(os.Getenv("TABLE_NAME")),
			Key: map[string]types.AttributeValue{
				"PK": &types.AttributeValueMemberS{Value: pkOrder},
				"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf(skOrder, orderID)},
			},
			ExpressionAttributeNames: expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			UpdateExpression: expr.Update(),
		})
		if err != nil {
			return attributeMap, fmt.Errorf("Couldn't update order with id: %s Because of this: %v\n", orderID, err)
		}else{
			err = attributevalue.UnmarshalMap(response.Attributes, &attributeMap)
			if err != nil {
				return attributeMap, fmt.Errorf("Couldn't unmarshall update response: %v\n", err)
			}
			fmt.Printf("Order with id %s was succesfuly updated with status %s", orderID, newStatus)
		}
		return attributeMap, err
	}
}