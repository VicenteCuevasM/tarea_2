package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"notificacion/models"
)

var client = &http.Client{}

const url = "https://sjwc0tz9e4.execute-api.us-east-2.amazonaws.com/Prod"

func SendNotification(order models.OrderMessage) error {

	var requestBody models.RequestBody
	requestBody.OrderID = order.OrderID
	requestBody.GroupID = "M6x!3F7l#1"
	requestBody.Products = order.Products
	requestBody.Customer = order.Customer

	//JSON encode
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	//Create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	//Send request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	//acepta 200 y 201, en el pdf decia 201 pero la request devuelve 200 como status code
	//y 201 como un campo del json
	if resp.StatusCode > 202 {
		fmt.Println("Error: ", resp.StatusCode)
		return err
	}
	defer resp.Body.Close()

	return nil
}
