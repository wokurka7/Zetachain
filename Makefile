.PHONY: build


PACKAGES=$(shell go list ./... | grep -v '/simulation')
VERSION := $(shell git describe --tags)
COMMIT := $(shell [ -z "${COMMIT_ID}" ] && git log -1 --format='%H' || echo ${COMMIT_ID} )
BUILDTIME := $(shell date -u +"%Y%m%d.%H%M%S" )
DOCKER ?= docker
DOCKER_BUF := $(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace bufbuild/buf
GOFLAGS:=""

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=zetacore \
	-X github.com/cosmos/cosmos-sdk/version.ServerName=zetacored \
	-X github.com/cosmos/cosmos-sdk/version.ClientName=zetaclientd \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
	-X github.com/zeta-chain/zetacore/common.Version=$(VERSION) \
	-X github.com/zeta-chain/zetacore/common.CommitHash=$(COMMIT) \
	-X github.com/zeta-chain/zetacore/common.BuildTime=$(BUILDTIME) \
	-X github.com/cosmos/cosmos-sdk/types.DBBackend=pebbledb

BUILD_FLAGS := -ldflags '$(ldflags)' -tags PRIVNET,pebbledb 
TEST_DIR?="./..."
TEST_BUILD_FLAGS := -tags PRIVNET,pebbledb 

clean: clean-binaries clean-dir

clean-binaries:
	@rm -rf ${GOBIN}/zetacored
	@rm -rf ${GOBIN}/zetaclientd

clean-dir:
	@rm -rf ~/.zetacored
	@rm -rf ~/.zetacore

all: install

test-coverage-exclude-core:
	@go test {TEST_BUILD_FLAGS} -v -coverprofile coverage.out $(go list ./... | grep -v /x/zetacore/)

test-coverage:
	@go test ${TEST_BUILD_FLAGS} -v -coverprofile coverage.out ${TEST_DIR}

coverage-report: test-coverage
	@go tool cover -html=cover.txt

test:
	@go test ${TEST_BUILD_FLAGS} ${TEST_DIR}

gosec:
	gosec  -exclude-dir=localnet ./...

install: go.sum
		@echo "--> Installing zetacored & zetaclientd"
		@go install -mod=readonly $(BUILD_FLAGS) ./cmd/zetacored
		@go install -mod=readonly $(BUILD_FLAGS) ./cmd/zetaclientd

install-zetaclient: go.sum
		@echo "--> Installing zetaclientd"
		@go install -mod=readonly $(BUILD_FLAGS) ./cmd/zetaclientd

# running with race detector on will be slow
install-zetaclient-race-test-only-build: go.sum
		@echo "--> Installing zetaclientd"
		@go install -race -mod=readonly $(BUILD_FLAGS) ./cmd/zetaclientd

install-zetacore: go.sum
		@echo "--> Installing zetacored"
		@go install -mod=readonly $(BUILD_FLAGS) ./cmd/zetacored

install-smoketest: go.sum
		@echo "--> Installing orchestrator"
		@go install -mod=readonly $(BUILD_FLAGS) ./contrib/localnet/orchestrator/smoketest

go.sum: go.mod
		@echo "--> Ensure dependencies have not been modified"
		GO111MODULE=on go mod verify

test-cctx:
	./standalone-network/cctx-creator.sh

init:
	./standalone-network/init.sh

run:
	./standalone-network/run.sh

chain-init: clean install-zetacore init
chain-run: clean install-zetacore init run

lint-pre:
	@test -z $(gofmt -l .)
	@GOFLAGS=$(GOFLAGS) go mod verify

lint: lint-pre
	@golangci-lint run

proto-go:
	@echo "--> Generating protobuf files"
	@ignite generate proto-go -y

###############################################################################
###                                Docker Images                             ###
###############################################################################

zetanode:
	@echo "Building zetanode"
	@docker build --no-cache -t zetanode -f ./Dockerfile .
.PHONY: zetanode

smoketest:
	@echo "--> Building smoketest image"
	$(DOCKER) build --no-cache -t orchestrator -f contrib/localnet/orchestrator/Dockerfile .
.PHONY: smoketest

start-smoketest:
	@echo "--> Starting smoketest"
	cd contrib/localnet/ && $(DOCKER) compose up -d

stop-smoketest:
	@echo "--> Stopping smoketest"
	cd contrib/localnet/ && $(DOCKER) compose down --remove-orphans


containerProtoVer=v0.7
containerProtoImage=tendermintdev/sdk-proto-gen:$(containerProtoVer)
containerProtoGen=cosmos-sdk-proto-gen-$(containerProtoVer)
containerProtoGenSwagger=cosmos-sdk-proto-gen-swagger-$(containerProtoVer)
containerProtoFmt=cosmos-sdk-proto-fmt-$(containerProtoVer)

proto-all: proto-format proto-lint proto-gen

proto-gen:
	@echo "Generating Protobuf files"
	$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace tendermintdev/sdk-proto-gen sh ./scripts/protocgen.sh

proto-format:
	@echo "Formatting Protobuf files"
	@if docker ps -a --format '{{.Names}}' | grep -Eq "^${containerProtoFmt}$$"; then docker start -a $(containerProtoFmt); else docker run --name $(containerProtoFmt) -v $(CURDIR):/workspace --workdir /workspace tendermintdev/docker-build-proto \
		find ./proto -name "*.proto" -exec clang-format -i {} \; ; fi

proto-lint:
	@$(DOCKER_BUF) lint --error-format=json

.PHONY: proto-all proto-gen proto-format proto-lint