# 📬 Notification Service

A lightweight, scalable notification service built with Go that consumes messages from RabbitMQ and sends notifications via email.

## ✨ Features

- 🐰 RabbitMQ integration for reliable message queuing
- 📧 Email notification delivery
- 🔄 Graceful shutdown handling
- 🛠️ Easy configuration via environment variables
- 🐳 Docker and Docker Compose support

## 🏗️ Architecture

The service follows a clean architecture approach:

- `cmd/app`: Entry point of the application
- `internal/app`: Application initialization and lifecycle management
- `internal/config`: Configuration management
- `internal/consumer`: RabbitMQ consumer implementation
- `internal/notification_service`: Notification service interface
- `internal/usecases`: Business logic implementation

## 🚀 Getting Started

### Prerequisites

- Docker and Docker Compose installed on your system

### Configuration

1. Create a `.env` file based on the provided `.env.example`:

```
MAIL_API_KEY="your-mail-api-key"
EMAIL_FROM_NAME="Your Service Name"
EMAIL_FROM="notifications@example.com"

RABBIT_URL="amqp://guest:guest@localhost:5672/"
RABBIT_QUEUE="notifications"
RABBIT_EXCHANGE="notifications"
RABBIT_ROUTING_KEY="email"
```

### 🐳 Running with Docker Compose

1. Clone the repository:

```bash
git clone https://github.com/idmaksim/notification-service.git
cd notification-service
```

2. Configure your environment variables in `.env` file

3. Build and start the services:

```bash
docker-compose up --build -d
```

This will start:

- The notification service application
- RabbitMQ with management console (accessible at http://localhost:15672)

4. To stop the services:

```bash
docker-compose down
```

## 📝 Message Format

To send a notification, publish a message to RabbitMQ with the following JSON format:

```json
{
  "target": "recipient@example.com",
  "subject": "Your notification subject",
  "text": "Your notification content"
}
```

## 📤 Sending a Test Message

You can send a test message to the broker using the provided script:

```bash
go run scripts/produce_message.go
```

This script will publish a test message to RabbitMQ using the settings from your `.env` file.

## 🔍 Monitoring

You can monitor RabbitMQ using the management console:

- URL: http://localhost:15672
- Username: guest
- Password: guest

## 🙏 Acknowledgements

- [Go](https://golang.org/)
- [RabbitMQ](https://www.rabbitmq.com/)
- [Docker](https://www.docker.com/)
