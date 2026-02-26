.PHONY: build test lint fmt vet clean run tidy

BINARY := server
BUILD_DIR := ./cmd
LDFLAGS := -s -w

build: ## Build the Go binary
	CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" -o $(BINARY) $(BUILD_DIR)

run: build ## Build and run
	./$(BINARY)

test: ## Run tests
	go test ./...

test-verbose: ## Run tests with verbose output
	go test -v ./...

test-cover: ## Run tests with coverage
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

lint: ## Run golangci-lint
	golangci-lint run ./...

fmt: ## Format code
	gofmt -w .

vet: ## Run go vet
	go vet ./...

tidy: ## Tidy go.mod
	go mod tidy

clean: ## Remove build artifacts
	rm -f $(BINARY) coverage.out
