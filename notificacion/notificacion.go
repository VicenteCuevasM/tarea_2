package main

import (
	"log"

	rabbit "notificacion/rabbitmq"
	"notificacion/service"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s\n", msg, err)
	}
}

func main() {

	// Initialize RabbitMQ
	err := rabbit.InitializeQueue("notificacion")
	if err != nil {
		failOnError(err, "Failed to initialze")
		return
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
		return
	}

	forever := make(chan bool)

	// Consume messages
	go func() {
		for d := range msgs {
			// Deserialize message
			orderMessage, err := rabbit.Deserialize(d.Body)
			if err != nil {
				failOnError(err, "Failed to deserialize")
			}

			// Send notification
			err = service.SendNotification(orderMessage)
			if err != nil {
				failOnError(err, "Failed to send notification")
			}

			log.Printf("Notification sent for order %s\n", orderMessage.OrderID)

		}
	}()

	log.Println(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

}
