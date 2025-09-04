export MYSQL_URL='mysql://root:secretPassword@tcp(localhost:3306)/elibrary'

# Migration commands
migrate-create:
	@ migrate create -ext sql -dir scripts/migrations -seq $(name)

migrate-up:
	@ migrate -database $(MYSQL_URL) -path scripts/migrations up

migrate-down:
	@ migrate -database $(MYSQL_URL) -path scripts/migrations down

migrate-force:
	@ migrate -database $(MYSQL_URL) -path scripts/migrations force $(version)

# Docker commands
docker-up:
	@ docker-compose up -d

docker-down:
	@ docker-compose down

docker-logs:
	@ docker-compose logs -f

# Application commands
run:
	@ go run cmd/main.go

build:
	@ go build -o bin/elibrary-backend cmd/main.go

test:
	@ go test ./...

# Development commands
dev: docker-up migrate-up run

install:
	@ go mod tidy
	@ go mod download

.PHONY: migrate-create migrate-up migrate-down migrate-force docker-up docker-down docker-logs run build test dev install
