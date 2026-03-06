APP_NAME=app

build:
	go build -o $(APP_NAME) cmd/app/main.go

run:
	go run cmd/app/main.go

tidy:
	go mod tidy

lint:
	golangci-lint run

tests:
	go test ./...

migrate-up:
	migrate -path migrations -database $(DB_URL) up

migrate-down:
	migrate -path migrations -database $(DB_URL) down

docker-build:
	docker compose build

docker-up:
	docker compose up -d

docker-down:
	docker compose down