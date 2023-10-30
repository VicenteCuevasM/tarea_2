package rabbitmq

import (
	"bytes"
	"encoding/json"
	"notificacion/models"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

var Conn *amqp.Connection
var Ch *amqp.Channel
var Q amqp.Queue

func InitializeQueue(nombre string) error {

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
		"tarea2", // name
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

	Q, err = Ch.QueueDeclare(
		nombre, // name
		false,  // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	if err != nil {
		return err
	}

	err = Ch.QueueBind(
		Q.Name,   // queue name
		"",       // routing key
		"tarea2", // exchange
		false,
		nil)
	if err != nil {
		return err
	}

	return nil
}

func Deserialize(msg []byte) (models.OrderMessage, error) {
	var orderMessage models.OrderMessage
	buf := bytes.NewBuffer(msg)
	dec := json.NewDecoder(buf)
	err := dec.Decode(&orderMessage)
	return orderMessage, err
}
