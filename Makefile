.PHONY: run test test-race cover lint fmt up down build

run:        ## Ejecuta la API localmente
	go run ./cmd/api

build:      ## Compila el binario
	go build -o bin/api ./cmd/api

test:       ## Tests unitarios
	go test ./...

test-race:  ## Tests con detector de condiciones de carrera
	go test -race -count=1 ./...

cover:      ## Reporte de cobertura
	go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

lint:       ## Linter (requiere golangci-lint)
	golangci-lint run ./...

fmt:
	gofmt -s -w .

up:         ## Levanta Postgres + API con Docker
	docker compose up -d --build

down:
	docker compose down
