.PHONY: build test run clean fmt lint

# Build the exporter binary
build:
	@echo "Building kubernetes-gateway-exporter..."
	@mkdir -p bin
	go build -o bin/exporter ./cmd/exporter

# Run tests
test:
	@echo "Running tests..."
	go test ./... -v

# Format Go code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint Go code
lint:
	@echo "Linting code..."
	go vet ./...

# Run the exporter locally
run: build
	@echo "Running exporter locally..."
	./bin/exporter

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
