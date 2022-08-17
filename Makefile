SHELL:=/bin/bash

define print =
	@echo "--- $(1)"
endef

define println =
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
	make lint test build

emerald: ./**/*.go
	$(benchmark_build_time)

iem: ./**/*.go
	$(benchmark_build_time)

build: emerald iem

test:
	go test ./lexer ./parser ./compiler/ ./object/ ./vm/ ./core/ -cover

lint:
	staticcheck ./...

install:
	go mod download && \
	go install honnef.co/go/tools/cmd/staticcheck@latest
