package main

import (
	"log"
	"os"

	"despacho/models"
	db "despacho/mongo"
	rabbit "despacho/rabbitmq"
	"despacho/utils"

	"github.com/joho/godotenv"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s\n", msg, err)
	}
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	mongo_string := "mongodb://" + os.Getenv("MV2_IP") + ":" + os.Getenv("MONGO_PORT")

	//Initialize MongoDB
	db.SetMongoString(mongo_string, "tarea2", "orders")
	err = db.InitMongoConnection()
	if err != nil {
		failOnError(err, "Failed to connect to MongoDB")
	}

	// Initialize RabbitMQ
	err = rabbit.InitializeQueue("despacho")
	if err != nil {
		failOnError(err, "Failed to initialze")
	}
	defer rabbit.Conn.Close()
	defer rabbit.Ch.Close()

	msgs, err := rabbit.Ch.Consume(
		rabbit.Q.Name, // queue
		"",            // consumer
		true,          // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	// Consume messages from RabbitMQ
	go func() {

		for d := range msgs {

			// Deserialize message
			orderMessage, err := rabbit.Deserialize(d.Body)
			if err != nil {
				failOnError(err, "Failed to deserialize message")
			}

			//Create Deliveries struct
			delivery := models.Delivery{
				ShippingAddress: models.ShippingAddress{
					Name:       orderMessage.Customer.Name,
					Lastname:   orderMessage.Customer.Lastname,
					Address1:   orderMessage.Customer.Location.Address1,
					Address2:   orderMessage.Customer.Location.Address2,
					City:       orderMessage.Customer.Location.City,
					State:      orderMessage.Customer.Location.State,
					PostalCode: orderMessage.Customer.Location.PostalCode,
					Country:    orderMessage.Customer.Location.Country,
					Phone:      orderMessage.Customer.Phone,
				},
				ShippingMethod: "UTFSM Delivery",
				TrackingNumber: utils.MakeTrackingNumber()}

			//Aqui no entiendo pq se usa un array si al final solo se envia un delivery
			//y me da miedo que se me haya pasado algo.
			deliveries := []models.Delivery{delivery}
			db.AddDeliveryToOrder(orderMessage.OrderID, deliveries)

			log.Printf("Delivery added to order: %s", orderMessage.OrderID)
		}

	}()

	log.Println(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
