.PHONY: build run clean install test help

# Variables
BINARY_NAME=tech-debt-collector
BINARY_DIR=bin
MAIN_PACKAGE=./cmd/tech-debt-collector
REPORT_FILE=report.json

help:
	@echo "Tech Debt Collector - Build & Run Commands"
	@echo ""
	@echo "Available targets:"
	@echo "  make build          - Build the binary"
	@echo "  make install        - Install binary to \$$GOPATH/bin"
	@echo "  make run            - Run analysis on current directory"
	@echo "  make run-test       - Run on test directory"
	@echo "  make clean          - Remove binary and reports"
	@echo "  make deps           - Download dependencies"
	@echo "  make fmt            - Format code"
	@echo "  make help           - Show this message"

build:
	@echo "🔨 Building $(BINARY_NAME)..."
	@mkdir -p $(BINARY_DIR)
	@go build -v -o $(BINARY_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "✅ Binary built: $(BINARY_DIR)/$(BINARY_NAME)"

install: build
	@echo "📦 Installing to \$$GOPATH/bin..."
	@go install $(MAIN_PACKAGE)
	@echo "✅ Installation complete"

run: build
	@echo "🚀 Running analysis on current directory..."
	@./$(BINARY_DIR)/$(BINARY_NAME) --path . --output $(REPORT_FILE) --verbose
	@echo "📊 Report saved to $(REPORT_FILE)"

run-llm: build
	@echo "🚀 Running analysis with LLM enrichment..."
	@./$(BINARY_DIR)/$(BINARY_NAME) --path . --output $(REPORT_FILE) --llm --verbose
	@echo "📊 Report saved to $(REPORT_FILE)"

run-text: build
	@echo "🚀 Running analysis (text output)..."
	@./$(BINARY_DIR)/$(BINARY_NAME) --path . --output report.txt --format text
	@echo "📊 Report saved to report.txt"
	@cat report.txt

run-test: build
	@echo "🚀 Running analysis on test directory..."
	@mkdir -p test_repo/src
	@echo 'package main\n\nfunc main() {\n  // TODO: Add error handling\n  // FIXME: Memory leak here\n}' > test_repo/src/test.go
	@./$(BINARY_DIR)/$(BINARY_NAME) --path test_repo --output test_report.json --verbose
	@echo "✅ Test complete - see test_report.json"

clean:
	@echo "🧹 Cleaning up..."
	@rm -rf $(BINARY_DIR)
	@rm -f $(REPORT_FILE) report.txt test_report.json
	@rm -rf test_repo
	@echo "✅ Clean complete"

deps:
	@echo "📥 Downloading dependencies..."
	@go mod download
	@go mod verify
	@echo "✅ Dependencies ready"

fmt:
	@echo "📝 Formatting code..."
	@go fmt ./...
	@echo "✅ Format complete"

lint:
	@echo "🔍 Linting code..."
	@golint ./...

test:
	@echo "🧪 Running tests..."
	@go test -v ./...

.DEFAULT_GOAL := help
