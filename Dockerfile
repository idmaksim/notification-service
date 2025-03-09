FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY .env ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/bin/service ./cmd/app


FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/bin/service ./service

COPY --from=builder /app/.env ./

EXPOSE 8080

CMD ["./service"]