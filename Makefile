BINARY_NAME=app.exe

all: test lint swagger build

start: build run

build: 
	@echo Building application...
	@go build -o $(BINARY_NAME) ./cmd/app
	@echo Application built successfully: $(BINARY_NAME)

run:
	@./$(BINARY_NAME)

test:
	@echo "Running tests..."
	@go test -v ./...

lint:
	@echo "Running linter..."
	@golangci-lint run


.DEFAULT_GOAL := build
