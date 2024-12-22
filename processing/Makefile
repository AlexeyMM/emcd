generate-protoc:
	@protoc \
		-I protocol/3rd-party \
		-I protocol/proto \
		--proto_path=. \
		--go_out=. \
		--go_opt=module='code.emcdtech.com/b2b/processing' \
		--go-grpc_out=. \
		--go-grpc_opt=module='code.emcdtech.com/b2b/processing' \
		./protocol/proto/*.proto

generate:
	go generate ./...

test: ### Run go test
	go test ./...

test_cover: ### Run go test with coverage estimate
	go test -v -coverprofile cover_test.out.tmp ./... && \
	cat cover_test.out.tmp | grep -v -E "mocks|mock\.go|pb\.go" > cover_test.out && \
	rm cover_test.out.tmp && \
	go tool cover -html cover_test.out -o cover_test.html && \
	open cover_test.html

generate-sqlc:
	# https://docs.sqlc.dev/en/latest/tutorials/getting-started-postgresql.html
	@sqlc generate -f sqlc.yaml

install-dev-deps:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.26.0
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.35.1
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/matryer/moq@v0.5.1

infra-up:
	docker compose -f deployments/local/docker-compose.deps.yml up -d

infra-stop:
	docker compose -f deployments/local/docker-compose.deps.yml stop

fmt:
	golangci-lint --config .golangci.yml run --fix

lint: fmt
	golangci-lint --config .golangci.yml run

