package adapters

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)
type MockDynamoDBClient struct{
	dynamodb.Client
}
func TestCreateOrderDynamo(t *testing.T){
	/*const pkOrder = "ORDER"
	const skOrder = "ID#%s"
	order := &dto.CreateOrderRequest{
		OrderID: "01",
		Item: "Item",
		UserID: "TestUser",
		Quantity: 2,
		TotalPrice: 4.3,
	} 
	
	t.Run("Add order correctly", func(t *testing.T){
		fakeDBClient := MockDynamoDBClient.Client.
	})*/
}