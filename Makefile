default:
	make lint test

test:
	go test ./lexer ./parser ./compiler/ ./vm/

lint:
	staticcheck ./...
