.PHONY: build run test lint clean tidy

APP_NAME := myapp
BUILD_DIR := bin

build:
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/myapp

run:
	go run ./cmd/myapp -c config.yaml

test:
	go test -race ./...

lint:
	golangci-lint run

tidy:
	go mod tidy

clean:
	rm -rf $(BUILD_DIR)
