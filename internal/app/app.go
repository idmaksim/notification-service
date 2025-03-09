package app

import (
	"github.com/idmaksim/notification-service/internal/config"
	"github.com/idmaksim/notification-service/internal/usecases/email"
)

type App struct {
	config  *config.Config
	service NotificationService
}

func NewApp() *App {
	return &App{
		config:  config.GetConfig(),
		service: email.NewMailService(),
	}
}

func (a *App) Run() error {
	err := a.service.Send(
		"dementevmaksim1657@gmail.com",
		"Hello",
		"Hello maksim",
	)
	return err
}
