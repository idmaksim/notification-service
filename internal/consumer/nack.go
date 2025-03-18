package consumer

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (c *Consumer) nack(delivery *amqp.Delivery) {
	err := delivery.Nack(false, true)
	if err != nil {
		log.Println("Error while acknowledging the message")
	}
}
