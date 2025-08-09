# TodoIng Backend Makefile
# Complete build and development tasks for proto-based architecture

# Variables
GO_VERSION := 1.23
PROTO_DIR := api/proto/v1
GENERATED_DIR := pkg/api/v1
BINARY_NAME := server
GRPC_BINARY_NAME := grpc-server
DOCKER_IMAGE := todoing-backend

# Proto generation tools
PROTOC := protoc
PROTOC_GEN_GO := protoc-gen-go
PROTOC_GEN_GO_GRPC := protoc-gen-go-grpc

# Default target
.PHONY: all
all: clean deps proto build

# Help target
.PHONY: help
help:
	@echo "TodoIng Backend Makefile Commands:"
	@echo ""
	@echo "Development:"
	@echo "  install-tools    Install protobuf tools and dependencies"
	@echo "  deps            Download Go dependencies"
	@echo "  proto           Generate Go code from proto files"
	@echo "  build           Build HTTP server binary"
	@echo "  build-grpc      Build gRPC server binary"
	@echo "  run             Run HTTP server locally"
	@echo "  run-grpc        Run gRPC server locally"
	@echo ""
	@echo "Testing:"
	@echo "  test            Run all tests"
	@echo "  test-unit       Run unit tests only"
	@echo "  test-convert    Run conversion layer tests"
	@echo "  test-coverage   Run tests with coverage report"
	@echo ""
	@echo "Quality:"
	@echo "  lint            Run Go linter"
	@echo "  fmt             Format Go code"
	@echo "  vet             Run Go vet"
	@echo ""
	@echo "Docker:"
	@echo "  docker-build         Build local Docker image"
	@echo "  docker-build-official Build official Docker image"
	@echo "  docker-run           Run local Docker container"
	@echo "  docker-run-official  Run official Docker container"
	@echo "  docker-pull          Pull official Docker image"
	@echo "  docker-stop          Stop and remove Docker containers"
	@echo "  docker-dev           Start Docker development environment"
	@echo ""
	@echo "Documentation:"
	@echo "  docs            Generate complete API documentation"
	@echo "  docs-build      Build documentation generator tool"
	@echo "  docs-clean      Clean generated documentation"
	@echo ""
	@echo "Utilities:"
	@echo "  clean           Clean build artifacts"
	@echo "  clean-proto     Clean generated proto files"

# Install protobuf tools
.PHONY: install-tools
install-tools:
	@echo "Installing protobuf tools..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "Protobuf tools installed successfully!"

# Download dependencies
.PHONY: deps
deps:
	@echo "Downloading Go dependencies..."
	go mod download
	go mod tidy
	@echo "Dependencies downloaded!"

# Generate Go code from proto files
.PHONY: proto
proto: install-tools
	@echo "Generating Go code from proto files..."
	@mkdir -p $(GENERATED_DIR)
	@for proto_file in $(PROTO_DIR)/*.proto; do \
		echo "Processing $$proto_file..."; \
		$(PROTOC) --go_out=$(GENERATED_DIR) --go_opt=paths=source_relative \
			--go-grpc_out=$(GENERATED_DIR) --go-grpc_opt=paths=source_relative \
			--proto_path=$(PROTO_DIR) $$proto_file; \
	done
	@echo "Proto code generation completed!"

# Build HTTP server
.PHONY: build
build: deps
	@echo "Building HTTP server..."
	go build -ldflags="-w -s" -o $(BINARY_NAME) ./cmd/api/main.go
	@echo "HTTP server built: $(BINARY_NAME)"

# Build gRPC server
.PHONY: build-grpc
build-grpc: deps proto
	@echo "Building gRPC server..."
	go build -ldflags="-w -s" -o $(GRPC_BINARY_NAME) ./cmd/grpc/main.go
	@echo "gRPC server built: $(GRPC_BINARY_NAME)"

# Run HTTP server locally
.PHONY: run
run: build
	@echo "Starting HTTP server..."
	./$(BINARY_NAME)

# Run gRPC server locally
.PHONY: run-grpc
run-grpc: build-grpc
	@echo "Starting gRPC server..."
	./$(GRPC_BINARY_NAME)

# Run all tests
.PHONY: test
test:
	@echo "Running all tests..."
	go test -v ./...

# Run unit tests only
.PHONY: test-unit
test-unit:
	@echo "Running unit tests..."
	go test -v -short ./...

# Run conversion layer tests
.PHONY: test-convert
test-convert:
	@echo "Running conversion layer tests..."
	go test -v ./internal/convert/...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Lint code
.PHONY: lint
lint:
	@echo "Running Go linter..."
	golangci-lint run

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting Go code..."
	go fmt ./...

# Vet code
.PHONY: vet
vet:
	@echo "Running Go vet..."
	go vet ./...

# Build Docker image
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):latest .
	@echo "Docker image built: $(DOCKER_IMAGE):latest"

# Build Docker image with official tag
.PHONY: docker-build-official
docker-build-official:
	@echo "Building official Docker image..."
	docker build -t axiu/todoing-go:latest -t axiu/todoing-go:dev .
	@echo "Official Docker image built: axiu/todoing-go:latest"

# Run Docker container
.PHONY: docker-run
docker-run:
	@echo "Running Docker container..."
	docker run -d --name todoing-go-backend \
		-p 5004:5004 \
		-e MONGO_URI="mongodb://localhost:27017/todoing" \
		-e JWT_SECRET="development_jwt_secret_key" \
		$(DOCKER_IMAGE):latest
	@echo "Docker container started: todoing-go-backend"

# Run official Docker image
.PHONY: docker-run-official
docker-run-official:
	@echo "Running official Docker image..."
	docker run -d --name todoing-go-official \
		-p 5004:5004 \
		-e MONGO_URI="mongodb://localhost:27017/todoing" \
		-e JWT_SECRET="development_jwt_secret_key" \
		axiu/todoing-go:latest
	@echo "Official Docker container started: todoing-go-official"

# Pull official Docker image
.PHONY: docker-pull
docker-pull:
	@echo "Pulling official Docker image..."
	docker pull axiu/todoing-go:latest
	@echo "Official Docker image pulled: axiu/todoing-go:latest"

# Stop and remove Docker containers
.PHONY: docker-stop
docker-stop:
	@echo "Stopping Docker containers..."
	-docker stop todoing-go-backend todoing-go-official 2>/dev/null || true
	-docker rm todoing-go-backend todoing-go-official 2>/dev/null || true
	@echo "Docker containers stopped and removed"

# Docker development environment
.PHONY: docker-dev
docker-dev:
	@echo "Starting Docker development environment..."
	cd .. && ./docker/deploy.sh golang up --build
	@echo "Docker development environment started"

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY_NAME)
	rm -f $(GRPC_BINARY_NAME)
	rm -f coverage.out
	rm -f coverage.html
	@echo "Clean completed!"

# Clean generated proto files
.PHONY: clean-proto
clean-proto:
	@echo "Cleaning generated proto files..."
	rm -rf $(GENERATED_DIR)
	@echo "Proto files cleaned!"

# Quick development setup
.PHONY: dev-setup
dev-setup: install-tools deps proto
	@echo "Development setup completed!"
	@echo "You can now run 'make build' or 'make build-grpc'"

# Verify installation
.PHONY: verify
verify:
	@echo "Verifying installation..."
	@echo "Go version: $(shell go version)"
	@echo "Protoc version: $(shell protoc --version || echo 'protoc not found')"
	@echo "protoc-gen-go: $(shell which protoc-gen-go || echo 'not found')"
	@echo "protoc-gen-go-grpc: $(shell which protoc-gen-go-grpc || echo 'not found')"
	@echo ""
	@echo "Proto files:"
	@find $(PROTO_DIR) -name "*.proto" | head -5
	@echo ""
	@echo "Generated files:"
	@find $(GENERATED_DIR) -name "*.go" 2>/dev/null | head -5 || echo "No generated files found"

# Documentation generation
.PHONY: docs-build
docs-build:
	@echo "Building documentation generator tool..."
	go build -o bin/generate-docs ./tools/generate_complete_api.go
	@echo "Documentation generator built: bin/generate-docs"

.PHONY: docs
docs: docs-build
	@echo "🚀 Generating complete API documentation..."
	./bin/generate-docs
	@echo "✅ Complete API documentation generated successfully!"
	@echo "📄 File location: docs/api_complete.json"
	@echo "📊 Documentation includes:"
	@echo "   - REST API endpoints with full specifications"
	@echo "   - gRPC API endpoints (simulated as HTTP)"
	@echo "   - Complete data model definitions"
	@echo "   - Request/response schemas with examples"
	@echo ""
	@echo "🔗 Access documentation:"
	@echo "   API Docs: http://localhost:5004/api-docs"
	@echo "   Swagger UI: http://localhost:5004/swagger/"
	@echo "   Documentation Center: http://localhost:5004/docs/"
	@echo ""
	@echo "💡 Usage:"
	@echo "   make run          # Start API server"
	@echo "   make docs         # Regenerate documentation"
	@echo "   make docs-clean   # Clean generated docs"

.PHONY: docs-clean
docs-clean:
	@echo "Cleaning generated documentation..."
	rm -f docs/api_complete.json
	rm -f bin/generate-docs
	@echo "Documentation files cleaned!"

.PHONY: docs-serve
docs-serve: docs run
	@echo "📖 API documentation is now available at:"
	@echo "   Complete API docs: http://localhost:5004/api-docs"
	@echo "   Swagger UI: http://localhost:5004/swagger/"
	@echo "   Documentation center: http://localhost:5004/docs/"

# Quick documentation update (for development)
.PHONY: docs-dev
docs-dev:
	@echo "🔄 Quick documentation update for development..."
	go run ./tools/generate_complete_api.go
	@echo "✅ Documentation updated! File: docs/api_complete.json"
