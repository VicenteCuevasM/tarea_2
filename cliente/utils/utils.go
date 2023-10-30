package utils

import (
	"cliente/models"
	"cliente/order"
	"encoding/json"
	"io"
	"os"
)

func LoadProducts() (models.OrderFile, error) {

	fileName := os.Args[1]

	var orderFile models.OrderFile
	file, err := os.Open(fileName)
	if err != nil {
		return orderFile, err
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return orderFile, err
	}

	err = json.Unmarshal(byteValue, &orderFile)
	if err != nil {
		return orderFile, err
	}

	return orderFile, nil

}

func StructToProtobuf(orderFile models.OrderFile) *order.OrderGRPCMessage {

	var orderGRPCMessage order.OrderGRPCMessage

	for _, product := range orderFile.Products {
		orderGRPCMessage.Products = append(orderGRPCMessage.Products, &order.Product{
			Title:       product.Title,
			Author:      product.Author,
			Genre:       product.Genre,
			Pages:       int32(product.Pages),
			Publication: product.Publication,
			Price:       product.Price,
			Quantity:    int32(product.Quantity),
		})
	}

	orderGRPCMessage.Customer = &order.Customer{
		Name:     orderFile.Customer.Name,
		Lastname: orderFile.Customer.Lastname,
		Email:    orderFile.Customer.Email,
		Location: &order.Location{
			Address1:   orderFile.Customer.Location.Address1,
			Address2:   orderFile.Customer.Location.Address2,
			City:       orderFile.Customer.Location.City,
			State:      orderFile.Customer.Location.State,
			PostalCode: orderFile.Customer.Location.PostalCode,
			Country:    orderFile.Customer.Location.Country,
		},
		Phone: orderFile.Customer.Phone,
	}

	return &orderGRPCMessage
}
