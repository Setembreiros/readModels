ENVIRONMENT=development

update:
	go mod tidy

build: update
	go build -o ./readModels cmd/main.go

run:
	go run ./cmd/main.go

run-dev:
	export ENVIRONMENT="${ENVIRONMENT}" && go run ./cmd/main.go

run-dev-windows: 
	set ENVIRONMENT=${ENVIRONMENT} && go run ./cmd/main.go

test:
	go generate -v ./internal/... && go test ./internal/...