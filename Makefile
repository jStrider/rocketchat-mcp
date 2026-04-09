BINARY     := rocketchat-mcp
GOTESTSUM  := go run gotest.tools/gotestsum@latest
GOLANGCI   := golangci-lint

.DEFAULT_GOAL := help
.PHONY: build test test-short coverage lint fmt vet check docker secrets clean tidy help

##@ Building

build: ## Build the binary
	CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o dist/$(BINARY) ./cmd/rocketchat-mcp/

##@ Testing

test: ## Run tests with coverage
	$(GOTESTSUM) --format pkgname-and-test-fails --format-hide-empty-pkg \
		-- -cover -race -covermode=atomic -coverprofile=coverage.out ./...

test-short: ## Run short tests only
	$(GOTESTSUM) --format pkgname-and-test-fails -- -short ./...

coverage: test ## Show coverage report
	go tool cover -func=coverage.out

##@ Code Quality

lint: ## Run golangci-lint
	$(GOLANGCI) run -c .golangci.yaml

fmt: ## Format code
	gofumpt -l -w .

vet: ## Run go vet
	go vet ./...

check: fmt vet lint test ## Run all checks (CI)

##@ Dependencies

tidy: ## Tidy go.mod
	go mod tidy

##@ Container

docker: ## Build Docker image
	docker build -t rocketchat-mcp:latest .

##@ Security

secrets: ## Scan for leaked secrets
	gitleaks detect --source .

##@ Misc

clean: ## Remove build artifacts
	rm -rf dist/ coverage.out

help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
