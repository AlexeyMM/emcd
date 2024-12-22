GO_INTERNAL_PKG = "code.emcdtech.com/emcd/service/email"
GO_SDK = "code.emcdtech.com/emcd/sdk"
GO_FILES = $(shell find . -type f -name '*.go' -not -path "*/mocks/*" -not -path "*.pb.go" -not -path "./protocol/email/*" | tr "\n" " ")
GO_FILES_GOLINES = $(shell find . -type f -name '*.go' -not -path "*/mocks/*" -not -path "*.pb.go" -not -path "*test.go" -not -path "./protocol/email/*" | tr "\n" " ")

install-tools: ### Install the tools necessary to work on the project
	go install github.com/vektra/mockery/v2@v2.40.1
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0
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

deps-up: ## Run all dependencies in docker for local development
	docker-compose -f ./deployments/local/docker-compose.deps.yml up -d

deps-down: ## Stop all docker dependencies used for local development
	docker-compose -f ./deployments/local/docker-compose.deps.yml down -v

docker-generate-protoc: ### Internal command that is called from the docker container to generate pb
	rm -r ./protocol/email/*** 2> /dev/null || (echo "dir ./protocol/email/ was empty"; exit 0)
	protoc \
		--proto_path=. \
		--go_out=. \
		--go_opt=module='code.emcdtech.com/emcd/service/email' \
		--go-grpc_out=. \
		--go-grpc_opt=module='code.emcdtech.com/emcd/service/email' \
		--validate_out='module=code.emcdtech.com/emcd/service/email,lang=go:.' \
		./protocol/proto/*.proto

generate: ### Run code generation for the project (protoc, go generate ./... and other tools if need)
	docker run --rm \
			--mount src="$(shell pwd)",target=/workspace,type=bind \
			--workdir /workspace \
			registry.emcdtech.com/docker-images/go-gen-proto:v1.0.1
	mockgen \
		-source protocol/email/email_grpc.pb.go \
		-destination protocol/email/email_grpc.mock.pb.go \
		-package email
	mockgen \
		-source protocol/email/email_provider_settings_grpc.pb.go \
		-destination protocol/email/email_provider_settings_grpc.mock.pb.go \
		-package email
	mockgen \
		-source protocol/email/email_template_grpc.pb.go \
		-destination protocol/email/email_template_grpc.mock.pb.go \
		-package email
	go generate ./...

local-grpc-healthcheck: ### Run grpccurl local healthcheck
	grpcurl -plaintext -d '{}' 127.0.0.1:8080 grpc.health.v1.Health/Check

local-pprof: ### Run pprof local
	go tool pprof goprofex http://localhost:9090/debug/pprof/allocs

local-metrics: ### Run curl metric local
	curl http://localhost:9090/metrics

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
