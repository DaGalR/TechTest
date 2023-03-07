package adapters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)
type HTTPPortClient interface{
	Do(req *http.Request) (*http.Response, error)
}
type HTTPPostUpdateOrder interface{
	CallUpdateOrdersService(string,string) error
}
type httpPostUpdateOrder struct{
	client HTTPPortClient
}


func NewHTTPClient(client HTTPPortClient) HTTPPostUpdateOrder{
	return &httpPostUpdateOrder{client}
}

func (h *httpPostUpdateOrder) CallUpdateOrdersService(orderID, new_status string) error {
	body := map[string]interface{}{"order_id": orderID, "new_status": new_status}
	jsonData, err := json.Marshal(body)
	fmt.Printf("RAW BODY %v\n", body)
	if err != nil{
		return fmt.Errorf("Error creating JSON body: %s", err.Error())
	}
	req, err := http.NewRequest(http.MethodPost,os.Getenv("ORDERS_URL"),bytes.NewBuffer(jsonData))
	if err != nil{
		return fmt.Errorf("Unable to generate request, see here: %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := h.client.Do(req)
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