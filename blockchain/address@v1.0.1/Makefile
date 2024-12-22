GO_INTERNAL_PKG = "code.emcdtech.com/emcd/blockchain/address"
GO_SDK = "code.emcdtech.com/emcd/sdk"
GO_FILES = $(shell find . -type f -name '*.go' -not -path "*/mocks/*" -not -path "*.pb.go" -not -path "./protocol/address/*" | tr "\n" " ")
GO_FILES_GOLINES = $(shell find . -type f -name '*.go' -not -path "*/mocks/*" -not -path "*.pb.go" -not -path "*test.go" -not -path "./protocol/address/*" | tr "\n" " ")

install-tools: ### Install the tools necessary to work on the project
	go install github.com/vektra/mockery/v2@v2.40.1
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

test_clean: ### Clean go test cache
	go clean -testcache

test: ### Run go test
	go test ./...

test_cover: ### Run go test with coverage estimate
	go test -v -coverprofile cover_test.out.tmp ./... && \
	cat cover_test.out.tmp | grep -v -E "mocks|mock\.go|pb\.go" > cover_test.out && \
	rm cover_test.out.tmp && \
	go tool cover -html cover_test.out -o cover_test.html && \
	open cover_test.html

linter: ### Run linter
	golangci-lint run ./...

deps-up: ## Run all dependencies in docker for local development
	docker-compose -f ./deployments/local/docker-compose.deps.yml up -d

deps-down: ## Stop all docker dependencies used for local development
	docker-compose -f ./deployments/local/docker-compose.deps.yml down -v

generate-proto: ### generate protoc file *.go
	protoc --go_out=protocol/ --go-grpc_out=protocol protocol/proto/address.proto

generate-mock: ### generate inline mock and mockery
	mockery
## inline generation
##	go generate ./...
## internal repository mock
#	mockery --name=AddressRepository --dir=internal/repository --filename=address.mock.go --output=internal/repository/mocks --outpkg=repository_mocks
## service mock
#	mockery --name=AddressService --dir=internal/service --filename=service.mock.go --output=internal/service/mocks  --outpkg=service_mocks
## grpc mock (handler)
#	mockery --name=AddressServiceClient --dir=protocol/address --filename=address_grpc.mock.pb.go --output=protocol/address/mocks --outpkg=proto_mocks
## external repository
#	mockery --name=AddressRepository --dir=repository --filename=address.mock.go --output=repository/mocks --outpkg=proto_mocks

local-grpc-healthcheck: ### Run grpccurl local healthcheck
	grpcurl -plaintext -d '{}' 127.0.0.1:8081 grpc.health.v1.Health/Check

local-pprof: ### Run pprof local
	go tool pprof goprofex http://localhost:8080/debug/pprof/allocs

local-metrics: ### Run curl metric local
	curl http://localhost:8080/metrics

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
