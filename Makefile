SHELL:=/bin/bash

PACKAGES := ./ ./parser/lexer ./parser/ast ./parser ./compiler/ ./object/ ./vm/ ./core/ ./heap/ ./bytecode/ ./internal/sandboxwire/

default:
	@make lint test-all build

emerald:
	@./scripts/build emerald

iem:
	@./scripts/build iem

emerald-sandbox-worker:
	@./scripts/build emerald-sandbox-worker

.PHONY: build emerald iem emerald-sandbox-worker test test-all ci-test coverage lint

build: emerald iem emerald-sandbox-worker

install:
	@make build && \
	cp ./emerald /usr/local/bin/emerald && \
	cp ./iem /usr/local/bin/iem && \
	cp ./emerald-sandbox-worker /usr/local/bin/emerald-sandbox-worker

test:
	@echo "Running test ${RUN}" && echo "" && \
	EM_TEST=1 go test $(PACKAGES) -run=${RUN}

test-all: coverage

ci-test:
	EM_TEST=1 go test $(PACKAGES)

coverage:
	@./scripts/generate_coverage.rb $(PACKAGES)

lint:
	@./scripts/ensure-command.sh staticcheck "go install honnef.co/go/tools/cmd/staticcheck@latest"
	@echo "Vetting ..." && go vet ./... && echo "Success ✓"
	@echo "Formatting ..." && go mod tidy && go fmt ./... && echo "Success ✓"
	@echo "Linting ..." && staticcheck ./... && echo "Success ✓"
