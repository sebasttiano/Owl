
# HELP
.PHONY: help tests test

help: ## This help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

tests: ## Run local tests
	go test -v -coverprofile=cover.txt `go list ./... | egrep -v 'proto|mock'`

lint: ## Run linters
	golangci-lint run -v ./...

install_lint: ## Get GLOLANGCI_LINT and install
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin $(GLOLANGCI_LINT_VERSION)

coverage: ## Check coverage
	go tool cover -func cover.txt

test-cover: tests coverage ## Run local tests with coverage checking

build: build-agent build-server  ## Build all components

build-agent: ## Build agent component
	go build -o ./cmd/clien/client ./cmd/client/

build-server: ## Build server component
	go build -o ./cmd/server/server ./cmd/server/

