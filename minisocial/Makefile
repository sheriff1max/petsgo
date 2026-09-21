.PHONY: run build docker-up docker-down migrate-up migrate-down

swagger:
	swag init -d ./cmd/api,./internal/handler,./internal/domain -g main.go

run:
	go run cmd/api/main.go

build:
	go build -o bin/minisocial cmd/api/main.go

docker-up:
	docker-compose up -d --build

docker-down:
	docker-compose down

migrate-up:
	psql -h localhost -U postgres -d minisocial -f migrations/001_init.sql

clean:
	rm -rf bin/