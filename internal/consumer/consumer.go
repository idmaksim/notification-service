package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/idmaksim/notification-service/internal/config"
	"github.com/idmaksim/notification-service/internal/notification_service"
	"github.com/idmaksim/notification-service/internal/usecases/email"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn           *amqp.Connection
	ch             *amqp.Channel
	notifyClosed   chan *amqp.Error
	done           chan struct{}
	reconnectDelay time.Duration
	cfg            *config.Config
	service        notificationService.NotificationService
}

func NewConsumer() (*Consumer, error) {
	cfg := config.GetConfig()
	consumer := &Consumer{
		cfg:            cfg,
		done:           make(chan struct{}),
		reconnectDelay: time.Second * 5,
		service:        email.NewMailService(),
	}

	var err error
	consumer.conn, err = amqp.Dial(cfg.RabbitUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %s", err)
	}

	consumer.ch, err = consumer.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %s", err)
	}

	consumer.notifyClosed = make(chan *amqp.Error)
	consumer.conn.NotifyClose(consumer.notifyClosed)

	return consumer, nil
}

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

func (c *Consumer) nack(delivery *amqp.Delivery) {
	err := delivery.Nack(false, true)
	if err != nil {
		log.Println("Error while acknowledging the message")
	}
}

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

func (c *Consumer) Stop() {
	close(c.done)
	if c.ch != nil {
		c.ch.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}
