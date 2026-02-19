ifneq (,$(wildcard .env))
	include .env
	export
endif

.PHONY: run dev build lint tidy migrate-up migrate-down migrate-create migrate-force migrate-version

run:
	go run ./cmd/main.go

dev:
	air

build:
	CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/auth-service ./cmd/main.go

lint:
	golangci-lint run ./...

tidy:
	go mod tidy

migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" down

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

migrate-force:
	migrate -path migrations -database "$(DB_URL)" force $(version)

migrate-version:
	migrate -path migrations -database "$(DB_URL)" version
