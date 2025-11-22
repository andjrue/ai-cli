APP_NAME=ai-cli

build:
	go build -o (bin/$(APP_NAME)) ./cmd

run:
	go run ./cmd main.go

tidy:
	go mod tidy

fmt:
	go fmt ./...

check-quality:
	make fmt
	make vet

vet:
	go vet ./...


