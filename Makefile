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
	@make lint test build

emerald: ./**/*.go
	$(benchmark_build_time)

iem: ./**/*.go
	$(benchmark_build_time)

build: emerald iem

install:
	@make build && \
	cp ./emerald /usr/local/bin/emerald && \
	cp ./iem /usr/local/bin/iem

test:
	@echo "Running test suite" && echo "" && \
 	go test ./parser/lexer ./parser ./compiler/ ./vm/ ./core/ --timeout=1s -coverprofile=./tmp/coverage.out && \
 	go tool cover -html=tmp/coverage.out -o tmp/coverage.html && echo ""

ci-test:
	go test ./parser/lexer ./parser ./compiler/ ./object/ ./vm/ ./core/

lint:
	@echo "Linting ..." && staticcheck ./... && echo ""
