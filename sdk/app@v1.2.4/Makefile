GO_INTERNAL_PKG = "code.emcdtech.com/emcd/sdk/app"
GO_SDK = "code.emcdtech.com/emcd/sdk"
GO_FILES = $(shell find . -type f -name '*.go' | tr "\n" " ")
GO_FILES_GOLINES = $(shell find . -type f -name '*.go' -not -path "*test.go" | tr "\n" " ")

install-tools: ### Install the tools necessary to work on the project
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2
	go install github.com/segmentio/golines@v0.11.0
	go install mvdan.cc/gofumpt@v0.5.0
	go install github.com/daixiang0/gci@v0.11.2

fmt: ### Code formatting
	golines --base-formatter=gofmt --ignore-generated -m 130 -w $(GO_FILES_GOLINES)
	gofumpt -w $(GO_FILES)
	gci write \
		--section Standard \
		--section Default \
		--section "Prefix($(GO_INTERNAL_PKG))" \
		--section "Prefix($(GO_SDK))" \
		$(GO_FILES) > /dev/null 2>&1

linter: ### Run linter
	golangci-lint run ./...

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
