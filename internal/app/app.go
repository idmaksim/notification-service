package app

import (
	"context"
	"fmt"
	"github.com/idmaksim/notification-service/internal/config"
	"github.com/idmaksim/notification-service/internal/consumer"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	config   *config.Config
	consumer *consumer.Consumer
}

func NewApp() *App {
	appConsumer, err := consumer.NewConsumer()
	if err != nil {
		panic(err)
	}
	return &App{
		config:   config.GetConfig(),
		consumer: appConsumer,
	}
}

func (a *App) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	if err := a.consumer.Start(ctx); err != nil {
		return fmt.Errorf("error starting consumers: %w", err)
	}

	log.Println("Notification service has started. Waiting for messages...")

	<-sigCh
	log.Println("Received stop signal. Завершение работы...")

	a.consumer.Stop()

	return nil
}
