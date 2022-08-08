CONTRACT=$(notdir $(shell pwd))
FILES_GO=$(patsubst %, go/$(CONTRACT)/%.go, $(FILES)) go/main.go
FILES_RS=$(patsubst %, src/%.rs, $(FILES))
FILES_TS=$(patsubst %, ts/$(CONTRACT)/%.ts, $(FILES)) ts/$(CONTRACT)/index.ts ts/$(CONTRACT)/tsconfig.json

all: build

schema-tool-install:
	@make --no-print-directory -C $(TOP_DIR) schema-tool-install

build: build-go build-rust build-ts

build-go: schema-tool-install $(FILES_GO)

$(FILES_GO): schema.yaml go/$(CONTRACT)/$(CONTRACT).go
	@echo "Generating Go files for contract $(CONTRACT)"
	schema -go
	golangci-lint run --fix

build-rust: schema-tool-install $(FILES_RS)

$(FILES_RS): schema.yaml src/$(CONTRACT).rs
	@echo "Generating Rust files for contract $(CONTRACT)"
	schema -rust

build-ts: schema-tool-install $(FILES_TS)

$(FILES_TS): schema.yaml ts/$(CONTRACT)/$(CONTRACT).ts
	@echo "Generating TypeScript files for contract $(CONTRACT)"
	schema -ts

clean: clean-go clean-rust clean-ts

clean-go:
	rm $(FILES_GO)

clean-rust:
	rm $(FILES_RS)

clean-ts:
	rm $(FILES_TS)
