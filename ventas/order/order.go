package order

import (
	"context"
	"log"
	"ventas/models"
	db "ventas/mongo"
	rabbit "ventas/rabbitmq"
)

type Server struct {
}

func (s *Server) CreateOrder(ctx context.Context, order *OrderGRPCMessage) (*OrderID, error) {
	log.Println("Order Received")

	orderStruct := convert(order)

	orderID, err := db.AddOrder(orderStruct)
	if err != nil {
		log.Printf("Error al insertar en la base de datos: %v", err)
		return nil, err
	}

	var rabbitMessage models.OrderRabbitMessage
	rabbitMessage.OrderID = orderID
	rabbitMessage.Products = orderStruct.Products
	rabbitMessage.Customer = orderStruct.Customer

	err = rabbit.PublishMessage(rabbitMessage)
	if err != nil {
		log.Printf("Error al publicar en RabbitMQ: %v", err)
		return nil, err
	}

	return &OrderID{Id: orderID}, nil
}

func convert(order *OrderGRPCMessage) models.OrderGRPCMessage {

	var products []models.Product
	for _, product := range order.Products {
		products = append(products, models.Product{
			Title:       product.Title,
			Author:      product.Author,
			Genre:       product.Genre,
			Pages:       product.Pages,
			Publication: product.Publication,
			Price:       product.Price,
			Quantity:    product.Quantity,
		})
	}

	var location models.Location
	location.Address1 = order.Customer.Location.Address1
	location.Address2 = order.Customer.Location.Address2
	location.City = order.Customer.Location.City
	location.State = order.Customer.Location.State
	location.PostalCode = order.Customer.Location.PostalCode
	location.Country = order.Customer.Location.Country

	var customer models.Customer
	customer.Name = order.Customer.Name
	customer.Lastname = order.Customer.Lastname
	customer.Email = order.Customer.Email
	customer.Location = location
	customer.Phone = order.Customer.Phone

	var orderStruct models.OrderGRPCMessage
	orderStruct.Products = products
	orderStruct.Customer = customer

	return orderStruct

}
