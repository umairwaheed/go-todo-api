.PHONY: help format lint test build run clean

# Default target
help:
	@echo "Available commands:"
	@echo "  format    - Format Go code using gofmt and goimports"
	@echo "  lint      - Run golangci-lint for code quality checks"
	@echo "  test      - Run all tests"
	@echo "  test-race - Run tests with race detection"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  build     - Build the application"
	@echo "  run       - Run the application"
	@echo "  clean     - Clean build artifacts and test files"
	@echo "  install-tools - Install required development tools"

# Format Go code
format:
	@echo "Formatting Go code..."
	@if command -v gofmt >/dev/null 2>&1; then \
		gofmt -s -w .; \
	else \
		echo "gofmt not found, using go fmt"; \
		go fmt ./...; \
	fi
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -w .; \
	else \
		echo "goimports not found, install with: go install golang.org/x/tools/cmd/goimports@latest"; \
	fi
	@echo "Code formatting complete!"

# Lint code using golangci-lint
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found, install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		echo "Running basic go vet instead..."; \
		go vet ./...; \
	fi

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with race detection
test-race:
	@echo "Running tests with race detection..."
	go test -race -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Build the application
build:
	@echo "Building application..."
	go build -o bin/todo-api main.go
	@echo "Binary created: bin/todo-api"

# Run the application
run:
	@echo "Running application..."
	go run main.go

# Clean build artifacts and test files
clean:
	@echo "Cleaning up..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	rm -f *.db
	rm -f test_*.db
	rm -f benchmark*.db
	@echo "Cleanup complete!"

# Install development tools
install-tools:
	@echo "Installing development tools..."
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Development tools installed!"

# Pre-commit hook: format, lint, and test
pre-commit: format lint test
	@echo "Pre-commit checks passed!"

# Development workflow: format, lint, test, and run
dev: format lint test run 