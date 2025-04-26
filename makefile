BINARY_NAME=GTS
BUILD_DIR=./GTS
CONFIG_PATH=./GTS/config/config.yaml
GO=go
GO_BUILD=$(GO) build
GO_TEST=$(GO) test
GO_CLEAN=$(GO) clean
GTS=$(BUILD_DIR)/$(BINARY_NAME)

.PHONY: all
all: build

.PHONY: build
build:
	mkdir -p $(BUILD_DIR)
	cp -r ./pkg/config $(BUILD_DIR)
	$(GO_BUILD) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/main

.PHONY: run
run: 
	$(GTS) --config $(CONFIG_PATH) ./...

.PHONY: test
test:
	$(GO_TEST) -v ./...

.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)
	rm -rf ./dist

.PHONY: project-run
project-run: 
	$(GTS) --config $(CONFIG_PATH) ./...

.PHONY: project-run-verbose
project-run-verbose: 
	$(GTS) --config $(CONFIG_PATH) --verbose ./...

.PHONY: run-skip-ts
run-skip-ts: 
	$(GTS) --config $(CONFIG_PATH) --skip-ts ./...

.PHONY: run-skip-scss
run-skip-scss: 
	$(GTS) --config $(CONFIG_PATH) --skip-scss ./...

.PHONY: run-skip-php
run-skip-php: 
	$(GTS) --config $(CONFIG_PATH) --skip-php ./...

.PHONY: php-to-html
php-to-html: 
	$(GTS) --config $(CONFIG_PATH) --skip-ts --skip-scss ./...

.PHONY: run-custom-output
run-custom-output: 
	$(GTS) --config $(CONFIG_PATH) -o $(BUILD_DIR)/output ./...

.PHONY: install
install: build
	mkdir -p $(GOPATH)/bin
	cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/

.PHONY: help
help:
	@echo "Make targets for the Extended Go Compiler:"
	@echo "  build               -   Build the compiler"
	@echo "  run                 -   Run the compiler with default config"
	@echo "  test                -   Run tests"
	@echo "  clean               -   Clean build artifacts"
	@echo "  project-run         -   run a project with the extended compiler"
	@echo "  project-run-verbose -   run with verbose output"
	@echo "  run-skip-ts         -   run without TypeScript compilation"
	@echo "  run-skip-scss       -   run without SCSS compilation"
	@echo "  run-skip-php        -   run without PHP to HTML conversion"
	@echo "  php-to-html         -   Only convert PHP to HTML"
	@echo "  run-custom-output   -   run with custom output location"
	@echo "  install             -   Install the compiler to GOPATH/bin"
