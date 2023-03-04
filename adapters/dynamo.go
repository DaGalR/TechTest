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

func CreateOrderDynamo(order *dto.CreateOrderRequest, client *dynamodb.Client) error {
	_, err := client.PutItem(context.Background(), &dynamodb.PutItemInput{
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

func CreatePaymentDynamo(payment *dto.CreatePaymentRequest, client *dynamodb.Client) error{
	order, err := GetOrderDynamo(payment.OrderID, client)
	if err != nil{
		return fmt.Errorf("Could not verify if order with id: " + payment.OrderID + " exists because of: " + err.Error())
	}
	if order == (&dto.CreateOrderRequest{}){
		return fmt.Errorf("There is no order with id: "+ payment.OrderID)
	}else{
		_, err = client.PutItem(context.Background(), &dynamodb.PutItemInput{
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
		fmt.Print("Updating order...\n")
		_, err := UpdateOrderStatusDynamo(payment.OrderID, "Complete", client)
		if err != nil{
			return err
		}
	}
	return nil
}

func GetOrderDynamo(orderID string, client *dynamodb.Client) (*dto.CreateOrderRequest, error) {
	order := dto.CreateOrderRequest{}
	fmt.Printf("LOOKING FOR ORDER WITH ID %s\n",orderID)
	data, err := client.GetItem(context.Background(), &dynamodb.GetItemInput{
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

func UpdateOrderStatusDynamo(orderID string, newStatus string, client *dynamodb.Client)(map[string]map[string]interface{}, error){
	update := expression.Set(expression.Name("status"), expression.Value(newStatus))
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	var response *dynamodb.UpdateItemOutput
	var attributeMap map[string]map[string]interface{}
	if err != nil{
		return attributeMap, fmt.Errorf("Couldn't build expression for order update: %v\n",err)
	}else{
		response, err = client.UpdateItem(context.Background(), &dynamodb.UpdateItemInput{
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
			fmt.Printf("Order with id %s was succesfuly updated with status complete", orderID)
		}
		return attributeMap, err
	}
}