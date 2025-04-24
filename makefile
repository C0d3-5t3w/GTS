BINARY_NAME=extendedgo
BUILD_DIR=./bin
CONFIG_PATH=./pkg/config/config.yaml
GO=go
GO_BUILD=$(GO) build
GO_TEST=$(GO) test
GO_CLEAN=$(GO) clean
EXT_GO=$(BUILD_DIR)/$(BINARY_NAME)

.PHONY: all
all: build

.PHONY: build
build:
	mkdir -p $(BUILD_DIR)
	$(GO_BUILD) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/main

.PHONY: run
run: build
	$(EXT_GO) --config $(CONFIG_PATH) ./...

.PHONY: test
test:
	$(GO_TEST) -v ./...

.PHONY: clean
clean:
	$(GO_CLEAN)
	rm -rf $(BUILD_DIR)
	rm -rf ./dist

.PHONY: project-build
project-build: build
	$(EXT_GO) --config $(CONFIG_PATH) ./...

.PHONY: project-build-verbose
project-build-verbose: build
	$(EXT_GO) --config $(CONFIG_PATH) --verbose ./...

.PHONY: build-skip-ts
build-skip-ts: build
	$(EXT_GO) --config $(CONFIG_PATH) --skip-ts ./...

.PHONY: build-skip-scss
build-skip-scss: build
	$(EXT_GO) --config $(CONFIG_PATH) --skip-scss ./...

.PHONY: build-custom-output
build-custom-output: build
	$(EXT_GO) --config $(CONFIG_PATH) -o $(BUILD_DIR)/output ./...

.PHONY: install
install: build
	mkdir -p $(GOPATH)/bin
	cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/

.PHONY: help
help:
	@echo "Make targets for the Extended Go Compiler:"
	@echo "  build                  Build the compiler"
	@echo "  run                    Run the compiler with default config"
	@echo "  test                   Run tests"
	@echo "  clean                  Clean build artifacts"
	@echo "  project-build          Build a project with the extended compiler"
	@echo "  project-build-verbose  Build with verbose output"
	@echo "  build-skip-ts          Build without TypeScript compilation"
	@echo "  build-skip-scss        Build without SCSS compilation"
	@echo "  build-custom-output    Build with custom output location"
	@echo "  install                Install the compiler to GOPATH/bin"
