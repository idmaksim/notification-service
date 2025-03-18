package consumer

import (
	"context"
	"log"
)

func (c *Consumer) handleCloseConnection(ctx *context.Context) {
	for {
		select {
		case <-c.done:
			return
		case err := <-c.notifyClosed:
			if err != nil {
				log.Printf("Connection with RabbitMQ is closed: %v", err)
				c.reconnect(*ctx)
			}
		}
	}
}
