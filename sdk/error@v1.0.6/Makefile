docker-generate-protoc: ### Internal command that is called from the docker container to generate pb
	rm -r ./proto/*.pb.go 2> /dev/null || (echo "not found files ./proto/*.pb.go"; exit 0)
	protoc \
		--proto_path=. \
		--go_out=. \
		--go_opt=module='code.emcdtech.com/emcd/sdk/error' \
		./proto/*.proto

generate: ### Run code generation for the project (protoc, go generate ./... and other tools if need)
	docker run --rm \
			--mount src="$(shell pwd)",target=/workspace,type=bind \
			--workdir /workspace \
			registry.emcdtech.com/docker-images/go-gen-proto:v0.0.1
	go generate ./...

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
