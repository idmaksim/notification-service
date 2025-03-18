package consumer

import (
	"context"
	"fmt"
)

func (c *Consumer) Start(ctx context.Context) error {
	if err := c.Setup(); err != nil {
		return err
	}

	msg, err := c.ch.Consume(
		c.cfg.RabbitQueue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error consume: %w", err)
	}

	go c.consume(ctx, msg)

	go c.handleCloseConnection(&ctx)

	return nil
}
