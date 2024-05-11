ENVIRONMENT=development

update:
	go mod tidy

build: update
	go build -o ./readModels cmd/main.go cmd/provider.go 

run:
	go run ./cmd/main.go cmd/provider.go

run-dev:
	export ENVIRONMENT="${ENVIRONMENT}" && go run ./cmd/main.go cmd/provider.go cmd/event_subscriptions.go

run-dev-windows: 
	set ENVIRONMENT=${ENVIRONMENT} && go run ./cmd/main.go cmd/provider.go cmd/event_subscriptions.go

test:
	go test ./...