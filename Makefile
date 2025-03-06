go-generate:
	@echo "Generating go-generate..."
	@go install github.com/vburenin/ifacemaker@v1.2.1
	@go install github.com/vektra/mockery/v2@v2.53.0
	@go generate ./internal/...

test:
	@echo "Running tests..."
	@go test ./internal/... -cover -race -short -count=1

lint:
	@echo "Running golangci-lint..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.6
	@golangci-lint run --fix ./...

run-server:
	@echo "Running server app..."
	@go run cmd/server/main.go

run-client:
	@echo "Running client app..."
	@go run cmd/client/main.go

docker-build:
	@echo "Docker build start..."
	@docker build -t wisdom-server .
	@docker build -t wisdom-client -f Dockerfile.client .
	@docker network create wisdom-net

docker-run-server:
	@echo "Running server docker app..."
	@docker run --rm -d --network=wisdom-net --name wisdom-server -p 9000:9000 wisdom-server

docker-run-client:
	@echo "Running client docker app..."
	@docker run --rm --network=wisdom-net wisdom-client




