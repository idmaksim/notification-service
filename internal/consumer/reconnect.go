package consumer

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (c *Consumer) reconnect(ctx context.Context) {
	for {
		select {
		case <-c.done:
			return
		case <-ctx.Done():
			return
		default:
			log.Printf("Trying to reconnect after %v...", c.reconnectDelay)
			time.Sleep(c.reconnectDelay)

			var err error
			c.conn, err = amqp.Dial(c.cfg.RabbitUrl)
			if err != nil {
				log.Printf("Error reconnecting to RabbitMQ: %v", err)
				continue
			}

			c.ch, err = c.conn.Channel()
			if err != nil {
				log.Printf("Error creating channel: %v", err)
				c.conn.Close()
				continue
			}

			c.notifyClosed = make(chan *amqp.Error)
			c.ch.NotifyClose(c.notifyClosed)

			if err := c.Start(ctx); err != nil {
				log.Printf("error starting consumer: %v", err)
				c.ch.Close()
				c.conn.Close()
				continue
			}

			log.Println("Successfully connected to RabbitMQ")
			return
		}
	}
}
