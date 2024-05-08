update:
	go mod tidy

build:
	go build -o readModels ./cmd

run:
	go run ./cmd/main.go

test:
	go test ./...