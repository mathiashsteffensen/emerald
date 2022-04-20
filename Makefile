default:
	make lint test build

build:
	go build -o ./tmp/emerald ./cmd/emerald/main.go && \
	go build -o ./tmp/iem ./cmd/iem/main.go

test:
	go test ./lexer ./parser ./compiler/ ./vm/ ./core/ -cover

lint:
	staticcheck ./...

install:
	go mod download && \
	go install honnef.co/go/tools/cmd/staticcheck@latest
