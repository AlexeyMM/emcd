docker-generate-protoc: ### Internal command that is called from the docker container to generate pb
	protoc \
		--proto_path=. \
		--go_out=. \
		--go_opt=module='code.emcdtech.com/b2b/b2b-endpoint' \
		--go-grpc_out=. \
		--go-grpc_opt=module='code.emcdtech.com/b2b/b2b-endpoint' \
		--validate_out='module=code.emcdtech.com/b2b/b2b-endpoint,lang=go:.' \
		./protocol/proto/*.proto

local-postgres-start:
	docker network create endpoint-network || true
	docker run --name endpoint-postgres --network endpoint-network -e POSTGRES_PASSWORD=mysecretpassword -p 5434:5432 -d postgres:16.3
	docker run --rm \
            --network endpoint-network \
            -v $(shell pwd)/migrations:/flyway/sql \
            flyway/flyway:latest \
            -url=jdbc:postgresql://endpoint-postgres:5432/postgres \
            -user=postgres \
            -password=mysecretpassword \
            migrate

local-postgres-stop:
	docker stop endpoint-postgres
	docker rm endpoint-postgres
	docker network rm endpoint-network

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

make_swagger:
	swag init --parseDependency

fmt:
	golangci-lint  --config .golangci-lint.yaml run --fix

lint: fmt
	golangci-lint  --config .golangci-lint.yaml run
