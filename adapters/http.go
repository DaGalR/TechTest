package adapters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func CallUpdateOrdersService(orderID string) error {
	body := map[string]interface{}{"order_id": orderID, "new_status": "Ready for shipping"}
	jsonData, err := json.Marshal(body)
	fmt.Printf("RAW BODY %v\n", body)
	if err != nil{
		return fmt.Errorf("Error creating JSON body: %s", err.Error())
	}
	resp, err := http.Post(os.Getenv("ORDERS_URL"), "application/json",bytes.NewBuffer(jsonData))
	if err != nil{
		return fmt.Errorf("Error sending POST request to URL: %s with error: %s", os.Getenv("ORDERS_URL"), err.Error())
	}
	var res map[string]interface{}

    err = json.NewDecoder(resp.Body).Decode(&res)
	fmt.Printf("RESULT FROM HTTP %v", res)
	if err != nil{
		return fmt.Errorf("Error reading response: %s", err.Error())
	}
	return nil
}