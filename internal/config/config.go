package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	MailApiKey       string `env:"MAIL_API_KEY"`
	EmailFromName    string `env:"EMAIL_FROM_NAME"`
	EmailFrom        string `env:"EMAIL_FROM"`
	RabbitUrl        string `env:"RABBIT_URL"`
	RabbitQueue      string `env:"RABBIT_QUEUE"`
	RabbitExchange   string `env:"RABBIT_EXCHANGE"`
	RabbitRoutingKey string `env:"RABBIT_ROUTING_KEY"`
}

var cfg *Config

func GetConfig() *Config {
	if cfg != nil {
		return cfg
	}

	return newConfig()
}

func newConfig() *Config {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	cfg = &Config{}

	err := env.Parse(cfg)
	if err != nil {
		panic(err)
	}

	return cfg

}
