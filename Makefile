GIT_COMMIT_SHA := $(shell git rev-list -1 HEAD)
BUILD_TAGS = rocksdb,builtin_static
BUILD_LD_FLAGS = "-X github.com/iotaledger/wasp/packages/wasp.VersionHash=$(GIT_COMMIT_SHA)"
GO_INSTALL_DIR=$(if $(GOBIN),$(GOBIN),$(if $(GOPATH),$(GOPATH/bin),~/go/bin))

#
# You can override these e.g. as
#     make test TEST_PKG=./packages/vm/core/testcore/ TEST_ARG="-v --run TestAccessNodes"
#
TEST_PKG=./...
TEST_ARG=

# List all the contracts that should be built before instalation
CONTRACTS=\
	corecontracts\
	dividend\
	donatewithfeedback\
	erc20\
	erc721\
	fairauction\
	fairroulette\
	helloworld\
	inccounter\
	schemacomment\
	testcore\
	testwasmlib\
	timestamp\
	tokenregistry\
	gascalibration/executiontime\
	gascalibration/memory\
	gascalibration/storage
WASM_CONTRACT_TARGETS=$(patsubst %, wasm-build-%, $(CONTRACTS))
WASM_CONTRACT_TARGETS_GO=$(patsubst %, wasm-build-%-go, $(CONTRACTS))
WASM_CONTRACT_TARGETS_RUST=$(patsubst %, wasm-build-%-rust, $(CONTRACTS))
WASM_CONTRACT_TARGETS_TS=$(patsubst %, wasm-build-%-ts, $(CONTRACTS))
WASM_CONTRACT_TARGETS_CLEAN=$(patsubst %, wasm-clean-%, $(CONTRACTS))

include tools/schema/shema_tool_files.mk
SCHEMA_TOOL_FILES_ABS=$(patsubst %, tools/schema/%, $(SCHEMA_TOOL_FILES))

all: build-lint

compile-solidity:
ifeq (, $(shell which solc))
	@echo "no solc found in PATH, evm contracts won't be compiled"
else
	cd packages/vm/core/evm/iscmagic && if ! git diff --quiet *.sol; then go generate; fi
	cd packages/evm/evmtest && if ! git diff --quiet *.sol; then go generate; fi
endif

build: compile-solidity
	go build -o . -tags $(BUILD_TAGS) -ldflags $(BUILD_LD_FLAGS) ./...

build-lint: build lint

test-full: install
	go test -tags $(BUILD_TAGS),runheavy ./... --timeout 60m --count 1 -failfast

test: install
	go test -tags $(BUILD_TAGS) $(TEST_PKG) --timeout 40m --count 1 -failfast $(TEST_ARG)

test-short:
	go test -tags $(BUILD_TAGS) --short --count 1 -failfast $(shell go list ./... | grep -v github.com/iotaledger/wasp/contracts/wasm | grep -v github.com/iotaledger/wasp/packages/vm/)

install: wasm-build compile-solidity
	go install -tags $(BUILD_TAGS) -ldflags $(BUILD_LD_FLAGS) ./...

lint:
	golangci-lint run --timeout 5m0s

gofumpt-list:
	gofumpt -l ./

docker-build:
	docker build \
		--build-arg BUILD_TAGS=${BUILD_TAGS} \
		--build-arg BUILD_LD_FLAGS='${BUILD_LD_FLAGS}' \
		.

schema-tool-install: $(GO_INSTALL_DIR)/schema
	@:	#to suppress "Nothing to be done" warning

$(GO_INSTALL_DIR)/schema: $(SCHEMA_TOOL_FILES_ABS)
	go install ./tools/schema

wasm-build: $(WASM_CONTRACT_TARGETS)

$(WASM_CONTRACT_TARGETS): wasm-build-%: wasm-build-%-go wasm-build-%-rust wasm-build-%-ts

$(WASM_CONTRACT_TARGETS_GO): wasm-build-%-go:
	@make --no-print-directory -C contracts/wasm/$* build-go

$(WASM_CONTRACT_TARGETS_RUST): wasm-build-%-rust:
	@make --no-print-directory -C contracts/wasm/$* build-rust

$(WASM_CONTRACT_TARGETS_TS): wasm-build-%-ts:
	@make --no-print-directory -C contracts/wasm/$* build-ts

wasm-clean: $(WASM_CONTRACT_TARGETS_CLEAN)

$(WASM_CONTRACT_TARGETS_CLEAN): wasm-clean-%:
	@make --no-print-directory -C contracts/wasm/$* clean

.PHONY: all compile-solidity build build-lint test-full test test-short install lint gofumpt-list docker-build schema-tool-install wasm-build $(WASM_CONTRACT_TARGETS) $(WASM_CONTRACT_TARGETS_GO) $(WASM_CONTRACT_TARGETS_RUST) $(WASM_CONTRACT_TARGETS_TS) wasm-clean $(WASM_CONTRACT_TARGETS_CLEAN)
