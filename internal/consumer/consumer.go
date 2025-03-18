package consumer

import (
	"fmt"

	"time"

	"github.com/idmaksim/notification-service/internal/config"
	"github.com/idmaksim/notification-service/internal/services"
	"github.com/idmaksim/notification-service/internal/usecases/email"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn           *amqp.Connection
	ch             *amqp.Channel
	notifyClosed   chan *amqp.Error
	done           chan struct{}
	reconnectDelay time.Duration
	cfg            *config.Config
	service        services.NotificationService
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
