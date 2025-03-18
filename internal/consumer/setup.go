package consumer

import (
	"fmt"
)

func (c *Consumer) Setup() error {
	err := c.ch.ExchangeDeclare(
		c.cfg.RabbitExchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error define exhange: %w", err)
	}

	_, err = c.ch.QueueDeclare(
		c.cfg.RabbitQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error define queue: %w", err)
	}

	err = c.ch.QueueBind(
		c.cfg.RabbitQueue,
		c.cfg.RabbitRoutingKey,
		c.cfg.RabbitExchange,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error queue bind: %w", err)
	}

	return nil
}
