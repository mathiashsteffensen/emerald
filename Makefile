default:
	make lint test build

emerald: ./**/*
	go build -o ./emerald ./cmd/emerald/main.go

iem: ./**/*
	go build -o ./iem ./cmd/iem/main.go

build: emerald iem

test:
	go test ./lexer ./parser ./compiler/ ./vm/ ./core/ -cover

lint:
	staticcheck ./...

install:
	go mod download && \
	go install honnef.co/go/tools/cmd/staticcheck@latest
