SHELL:=/bin/bash

default:
	@make lint test-all build

emerald:
	@./scripts/build emerald

iem:
	@./scripts/build iem

.PHONY: build emerald iem

build: emerald iem

install:
	@make build && \
	cp ./emerald /usr/local/bin/emerald && \
	cp ./iem /usr/local/bin/iem

test:
	@echo "Running test ${RUN}" && echo "" && \
	EM_TEST=1 go test ./parser/lexer ./parser ./compiler/ ./vm/ ./core/ -run=${RUN}

test-all:
	@echo "Running test suite" && echo "" && \
 	EM_TEST=1 go test ./parser/lexer ./parser ./compiler/ ./vm/ ./core/ --timeout=1s -coverprofile=./tmp/coverage.out && \
 	go tool cover -html=tmp/coverage.out -o tmp/coverage.html && echo ""

ci-test:
	EM_TEST=1 go test ./parser/lexer ./parser ./compiler/ ./object/ ./vm/ ./core/

lint:
	@echo "Linting ..." && staticcheck ./... && echo ""
