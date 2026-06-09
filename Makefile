.PHONY: build test run clean fmt fmt-check vet lint staticcheck tidy-check helm-lint docker-build verify

# Build the exporter binary
build:
	@echo "Building kubernetes-gateway-exporter..."
	@mkdir -p bin
	go build -o bin/exporter ./cmd/exporter

# Run tests
test:
	@echo "Running tests..."
	go test ./... -v

# Run tests with race detector
test-race:
	@echo "Running tests with race detector..."
	go test -race ./... -v

# Format Go code (rewrites files)
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Check formatting without rewriting — fails if any files are unformatted
fmt-check:
	@echo "Checking formatting..."
	@UNFORMATTED=$$(gofmt -l .); \
	if [ -n "$$UNFORMATTED" ]; then \
		echo "The following files are not gofmt-formatted:"; \
		echo "$$UNFORMATTED"; \
		exit 1; \
	fi
	@echo "All files are properly formatted."

# Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

# Run staticcheck (fetched on demand via go run — no manual install needed)
staticcheck:
	@echo "Running staticcheck..."
	go run honnef.co/go/tools/cmd/staticcheck@2025.1.1 ./...

# Lint: vet + staticcheck
lint:
	@echo "Linting code..."
	@$(MAKE) vet
	@$(MAKE) staticcheck

# Check that go.mod/go.sum are tidy. Compares each file before/after `go mod tidy`
# (not against git), so it works with uncommitted changes too.
tidy-check:
	@echo "Checking go mod tidy..."
	@cp go.mod go.mod.bak && cp go.sum go.sum.bak && go mod tidy; \
	if cmp -s go.mod go.mod.bak && cmp -s go.sum go.sum.bak; then \
		ec=0; echo "go.mod/go.sum are tidy."; \
	else \
		ec=1; echo "go.mod/go.sum are not tidy — run 'go mod tidy' and commit the result."; \
	fi; \
	mv go.mod.bak go.mod; mv go.sum.bak go.sum; exit $$ec

# Lint Helm chart
helm-lint:
	@echo "Linting Helm chart..."
	helm lint deploy/chart/kubernetes-gateway-exporter

# Build Docker image locally (no push)
docker-build:
	@echo "Building Docker image..."
	docker build -t kubernetes-gateway-exporter:dev .

# Run the exporter locally
run: build
	@echo "Running exporter locally..."
	./bin/exporter

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/

# Full local verification — mirrors what CI runs
verify: fmt-check vet staticcheck tidy-check test
	@echo ""
	@echo "All checks passed: fmt-check, vet, staticcheck, tidy-check, test."
