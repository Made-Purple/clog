.PHONY: build test test-coverage install clean

VERSION ?= dev

build:
	go build -ldflags "-X github.com/made-purple/clog/internal/command.Version=$(VERSION)" -o clog ./cmd/clog

test:
	go test ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

install:
	go install ./cmd/clog

clean:
	rm -f clog coverage.out
