APP_NAME=app
DB_URL=postgres://postgres:postgres@localhost:3333/tasks_db?sslmode=disable

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

database:
	PGPASSWORD=postgres psql -h localhost -p 3333 -U postgres -d tasks_db