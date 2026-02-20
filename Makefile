.PHONY: build run dev test docker-build docker-up docker-down clean

APP_NAME := goauthly-fiber-example
BUILD_DIR := ./bin

build:
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/server/

run: build
	$(BUILD_DIR)/$(APP_NAME)

dev:
	bash scripts/dev.sh

test:
	go test ./...

docker-build:
	docker build -t $(APP_NAME):latest .

docker-up:
	docker compose up --build -d

docker-down:
	docker compose down

clean:
	rm -rf $(BUILD_DIR)
