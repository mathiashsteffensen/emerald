default:
	make lint test build

build:
	go build -o ./tmp/emerald ./cmd/main.go && \
	go build -o ./tmp/iem ./main.go

test:
	go test ./lexer ./parser ./compiler/ ./vm/

lint:
	staticcheck ./...
