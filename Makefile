.DEFAULT_GOAL := help

check: test lint vet ## Runs all tests

test: ## Run the unit tests
	go test -cover -v $(shell go list ./... | grep -v /vendor/)

lint: ## Lint all files
	go list ./... | grep -v /vendor/ | xargs -L1 /Users/jbleduig/go/bin/golint -set_exit_status

vet: ## Run the vet tool
	go vet $(shell go list ./... | grep -v /vendor/)

clean: ## Clean up build artifacts
	# go clean
	rm -f budget2sheets budget2sheets.zip

build: clean test ## Build the executable
	GOOS=linux GOARCH=amd64 go build -o budget2sheets ./cmd/budget2sheets

zip: build ## Zip the executable so that it can be uploaded to AWS Lambda
	zip budget2sheets.zip budget2sheets

help: ## Display this help message
	@cat $(MAKEFILE_LIST) | grep -e "^[a-zA-Z_\-]*: *.*## *" | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.SILENT: zip build test lint vet clean help
