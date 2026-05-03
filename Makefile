SHELL:=/bin/bash

PACKAGES := ./parser/lexer ./parser ./compiler/ ./object/ ./vm/ ./core/ ./heap/ ./bytecode/

default:
	@make lint test-all build

emerald:
	@./scripts/build emerald

iem:
	@./scripts/build iem

.PHONY: build emerald iem test test-all ci-test coverage lint

build: emerald iem

install:
	@make build && \
	cp ./emerald /usr/local/bin/emerald && \
	cp ./iem /usr/local/bin/iem

test:
	@echo "Running test ${RUN}" && echo "" && \
	EM_TEST=1 go test $(PACKAGES) -run=${RUN}

test-all: coverage

ci-test:
	EM_TEST=1 go test $(PACKAGES)

coverage:
	@./scripts/generate_coverage.rb $(PACKAGES)

lint:
	@echo "Linting ..." && staticcheck ./... && echo ""
