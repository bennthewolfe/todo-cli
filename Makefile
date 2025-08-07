# Makefile for todo-cli project

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=todo
BINARY_WINDOWS=$(BINARY_NAME).exe

# Test parameters
TEST_PACKAGES=./...
COVERAGE_FILE=coverage.out

.PHONY: all build clean test coverage bench bench-verbose lint help

# Default target
all: test build

# Build the application
build:
	$(GOBUILD) -o $(BINARY_WINDOWS) -v .

# Clean build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_WINDOWS)
	rm -f $(COVERAGE_FILE)

# Run all tests
test:
	$(GOTEST) -v $(TEST_PACKAGES)

# Run tests with coverage
coverage:
	$(GOTEST) -v -coverprofile=$(COVERAGE_FILE) $(TEST_PACKAGES)
	$(GOCMD) tool cover -html=$(COVERAGE_FILE) -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run unit tests only (exclude integration tests)
test-unit:
	$(GOTEST) -v -short $(TEST_PACKAGES)

# Run integration tests only
test-integration:
	$(GOTEST) -v -run TestCLI $(TEST_PACKAGES)

# Run benchmark tests
bench:
	$(GOTEST) -bench=Benchmark -benchmem -run=^$

# Run benchmark tests with all tests (verbose)
bench-verbose:
	$(GOTEST) -bench=Benchmark -benchmem -v

# Run tests with race condition detection
test-race:
	$(GOTEST) -v -race $(TEST_PACKAGES)

# Lint the code (requires golangci-lint to be installed)
lint:
	golangci-lint run

# Format the code
fmt:
	$(GOCMD) fmt $(TEST_PACKAGES)

# Vet the code
vet:
	$(GOCMD) vet $(TEST_PACKAGES)

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Run all quality checks
check: fmt vet lint test

# Install the application
install:
	$(GOCMD) install .

# Show help
help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  clean         - Clean build files"
	@echo "  test          - Run all tests"
	@echo "  test-unit     - Run unit tests only"
	@echo "  test-integration - Run integration tests only"
	@echo "  coverage      - Run tests with coverage report"
	@echo "  bench         - Run benchmark tests only"
	@echo "  bench-verbose - Run benchmark tests with all tests"
	@echo "  test-race     - Run tests with race condition detection"
	@echo "  lint          - Lint the code"
	@echo "  fmt           - Format the code"
	@echo "  vet           - Vet the code"
	@echo "  check         - Run all quality checks"
	@echo "  deps          - Download dependencies"
	@echo "  install       - Install the application"
	@echo "  help          - Show this help"
