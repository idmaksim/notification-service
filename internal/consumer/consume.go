package consumer

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (c *Consumer) consume(ctx context.Context, deliveries <-chan amqp.Delivery) {
	for {
		select {
		case <-ctx.Done():
			return

		case delivery, ok := <-deliveries:
			log.Println("Message received")
			if !ok {
				log.Println("Channel is closed")
				return
			}

			var msg Message
			if err := json.Unmarshal(delivery.Body, &msg); err != nil {
				log.Println("Error while parsing the message")
				c.nack(&delivery)
				continue
			}

			if err := c.service.Send(msg.Target, msg.Subject, msg.Text); err != nil {
				log.Printf("error while sending notification: %v", err)
				c.nack(&delivery)
				continue
			}

			if err := delivery.Ack(false); err != nil {
				log.Printf("error accepting message: %v", err)
			}
		}
	}
}
