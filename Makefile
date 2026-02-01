.PHONY: build test lint clean help all

BINARY_NAME=kin-core-example
BUILD_DIR=bin

all: build test lint

help: ## Display this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the example application
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./example/main.go

test: ## Run tests
	go test -v ./...

lint: ## Run golangci-lint
	go tool golangci-lint run ./...

clean: ## Remove build artifacts
	rm -rf $(BUILD_DIR)
	go clean
