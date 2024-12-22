generate-protoc:
	@protoc \
		-I protocol/3rd-party \
		--proto_path=. \
		--go_out=. \
		--go_opt=module='code.emcdtech.com/b2b/swap' \
		--go-grpc_out=. \
		--go-grpc_opt=module='code.emcdtech.com/b2b/swap' \
		./protocol/proto/*.proto

local-postgres-start:
	docker network create mynetwork
	docker run --name my-postgres --network mynetwork -e POSTGRES_PASSWORD=mysecretpassword  -p 5432:5432 -d postgres
	docker run --rm \
            --network mynetwork \
            -v $(shell pwd)/migrations:/flyway/sql \
            flyway/flyway:latest \
            -url=jdbc:postgresql://my-postgres:5432/postgres \
            -user=postgres \
            -password=mysecretpassword \
            migrate

local-postgres-stop:
	docker stop my-postgres
	docker rm my-postgres
	docker network rm mynetwork


install-tools: ### Install the tools necessary to work on the project
	go install github.com/golang/mock/mockgen@v1.6.0
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2
	go install github.com/segmentio/golines@v0.11.0
	go install mvdan.cc/gofumpt@v0.5.0
	go install github.com/daixiang0/gci@v0.11.2


local-grpc-healthcheck: ### Run grpccurl local healthcheck
	grpcurl -plaintext -d '{}' 127.0.0.1:9090 grpc.health.v1.Health/Check

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

generate-mock: ### generate inline mock and mockery
	mockery

linter:
	golangci-lint run --config ./.golangci-lint.yaml
