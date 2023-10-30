package rabbitmq

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"time"
	"ventas/models"

	amqp "github.com/rabbitmq/amqp091-go"
)

var Conn *amqp.Connection
var Ch *amqp.Channel

func InitializeExchange(nombre string) error {

	rabbit_url := "amqp://guest:guest@" + os.Getenv("MV3_IP") + ":5672/"

	Conn, err := amqp.Dial(rabbit_url)
	if err != nil {
		return err
	}

	Ch, err = Conn.Channel()
	if err != nil {
		return err
	}

	err = Ch.ExchangeDeclare(
		nombre,   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return err
	}

	return nil
}

func PublishMessage(order models.OrderRabbitMessage) error {
	body, err := serialize(order)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = Ch.PublishWithContext(ctx,
		"tarea2", // exchange,
		"",       // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	if err != nil {
		return err
	}

	return nil
}

func serialize(order models.OrderRabbitMessage) ([]byte, error) {
	var b bytes.Buffer
	encoder := json.NewEncoder(&b)
	err := encoder.Encode(order)
	return b.Bytes(), err
}
