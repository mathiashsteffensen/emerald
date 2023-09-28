SHELL:=/bin/bash

define print
	@echo "--- $(1)"
endef

define println
	$(call print,$(1))
	@echo ""
endef

define get_timestamp
	$(shell date +%s%N)
endef

define benchmark_build_time =
	$(eval START:=$(get_timestamp))
	@echo ""
	$(call print, "Building $@ executable...")
	@go build -o ./$@ ./cmd/$@/main.go
	$(call println, "Completed building $@ executable in $$(((($(get_timestamp) - $(START))/1000)))μs")
endef

default:
	@make lint test-all build

emerald:
	$(benchmark_build_time)

iem: ./**/*.go ./parser/**/*.go ./cmd/iem/main.go
	$(benchmark_build_time)

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
