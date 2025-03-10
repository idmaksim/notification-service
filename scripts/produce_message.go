package main

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	//target = "irina_8823@mail.ru"
	target       = "dementevmaksim1657@gmail.com"
	subject      = "Subject"
	messageCount = 10
	text         = "Hello worl"

	rabbitMQURL = "amqp://guest:guest@localhost:5672/"
)

type NotificationMessage struct {
	Target  string `json:"target"`
	Subject string `json:"subject"`
	Text    string `json:"text"`
}

func main() {
	exchange := "notifications"
	routingKey := "email"

	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("Error connecting to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error creating channel: %v", err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Error  exchange defining: %v", err)
	}

	msg := NotificationMessage{
		Target:  target,
		Subject: subject,
		Text:    text,
	}

	body, err := json.Marshal(msg)
	if err != nil {
		log.Fatalf("Erro serialization the message: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	wg := sync.WaitGroup{}

	for i := 0; i < messageCount; i++ {
		wg.Add(1)
		go func() {
			wg.Done()
			err = ch.PublishWithContext(
				ctx,
				exchange,
				routingKey,
				false,
				false,
				amqp.Publishing{
					ContentType:  "application/json",
					DeliveryMode: amqp.Persistent,
					Body:         body,
				},
			)
			if err != nil {
				log.Fatalf("Error sending a message: %v", err)
			}
			log.Println("The message sent successfully")
		}()
	}

	wg.Wait()

}
