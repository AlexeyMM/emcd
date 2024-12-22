.PHONY: help # Generate list of targets with descriptions
help:
	@grep '^.PHONY: .* #' Makefile | sed 's/\.PHONY: \(.*\) # \(.*\)/\1 \2/' | expand -t20

.PHONY: install-tools # - Install necessary tools
install-tools: ### Install the tools necessary to work on the project
	go install github.com/vektra/mockery/v2@v2.45.0

test_clean: ### Clean go test cache
	go clean -testcache

test: ### Run go test
	go test ./...

test_cover: ### Run go test with coverage estimate
	go test -v -coverprofile cover_test.out.tmp ./... && \
	cat cover_test.out.tmp | grep -v -E "mock\.go|pb\.go" > cover_test.out && \
	rm cover_test.out.tmp && \
	go tool cover -html cover_test.out -o cover_test.html && \
	open cover_test.html

linter: ### Run linter
	golangci-lint run ./...

.PHONY: generate # - Generate protobuf and mocks
generate-proto:
	protoc --experimental_allow_proto3_optional --go_out=protocol/ --go-grpc_out=protocol/ protocol/accounting/accounting.proto
	protoc --go_out=protocol/ --go-grpc_out=protocol/ protocol/referral/referral_accounting.proto
	protoc --go_out=protocol/ --go-grpc_out=protocol/ protocol/user_account/user_account.proto

generate-mock: ### generate inline mock and mockery
	mockery

##	grpc mocks
#	mockery --name=AccountingServiceClient --dir=protocol/accounting/ --filename=accounting_grpc.mock.pb.go --output=protocol/accounting/mocks
#	mockery --name=AccountingReferralServiceClient --dir=protocol/referral/ --filename=referral_grpc.mock.pb.go --output=protocol/referral/mocks
#	mockery --name=UserAccountServiceClient --structname=UserAccountServiceClientRepo --with-expecter --dir=protocol/user_account/ --filename=user_account_grpc.mock.pb.go --output=protocol/user_account/mocks
##	external repository mocks
#	mockery --name=AccountingRepository --structname=MockAccountingRepository --with-expecter --dir=./repository --outpkg=repository_mock --output ./repository/repository_mock --filename balance_mock.go
#	mockery --name=HistoryRepository --structname=MockHistoryRepository --with-expecter --dir=./repository --outpkg=repository_mock --output ./repository/repository_mock --filename history_mock.go
#	mockery --name=UserAccountRepository --structname=MockUserAccountRepository --with-expecter --dir=./repository --outpkg=repository_mock --output ./repository/repository_mock --filename user_account_mock.go
##	service mocks
#	mockery --name=Balance --structname=MockBalance --dir=./internal/service --outpkg=mock --output ./internal/service/mock --filename balance_mock.go
#	mockery --name=UserAccountService --structname=MockUserAccountService --with-expecter --dir=./internal/service --outpkg=mock --output ./internal/service/mock --filename user_account_mock.go
##	db repository mocks
#	mockery --name=IncomesHistory --structname=MockIncomesHistory --with-expecter --dir=./internal/repository --outpkg=mock --output ./internal/repository/mock --filename income_history_mock.go
#	mockery --name=PayoutsHistory --structname=MockPayoutsHistory --with-expecter --dir=./internal/repository --outpkg=mock --output ./internal/repository/mock --filename payout_history_mock.go
#	mockery --name=WalletsHistory --structname=MockWalletsHistory --with-expecter --dir=./internal/repository --outpkg=mock --output ./internal/repository/mock --filename wallet_history_mock.go
#	mockery --name=ReferralStatistic --structname=MockReferralStatistic --with-expecter --dir=./internal/repository --outpkg=mock --output ./internal/repository/mock --filename referral_statistic_mock.go
#	mockery --name=Balance --structname=MockBalance --with-expecter --dir=./internal/repository --outpkg=mock --output ./internal/repository/mock --filename balance_mock.go
#	mockery --name=UserAccountRepo --structname=MockUserAccountRepo --with-expecter --dir=./internal/repository --outpkg=mock --output ./internal/repository/mock --filename user_account_mock.go


.DEFAULT_GOAL := help
