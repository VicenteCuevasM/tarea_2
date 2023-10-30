package main

import (
	"log"

	db "inventario/mongo"
	rabbit "inventario/rabbitmq"

	"os"

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
	db.SetMongoString(mongo_string, "tarea2", "products")
	err = db.InitMongoConnection()
	if err != nil {
		failOnError(err, "Failed to connect to MongoDB")
	}

	// Initialize RabbitMQ
	err = rabbit.InitializeQueue("inventario")
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
	if err != nil {
		failOnError(err, "Failed to register a consumer")
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			// Deserialize message
			orderMessage, err := rabbit.Deserialize(d.Body)
			if err != nil {
				failOnError(err, "Failed to deserialize message")
			}
			// Process order
			err = db.UpdateMultipleStock(orderMessage)
			if err != nil {
				failOnError(err, "Failed to update stock")
			}

			log.Printf("Order processed: %s", orderMessage.OrderID)
		}
	}()

	log.Println(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

}
