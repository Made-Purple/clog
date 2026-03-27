.PHONY: build test install clean

VERSION ?= dev

build:
	go build -ldflags "-X github.com/made-purple/clog/internal/command.Version=$(VERSION)" -o clog ./cmd/clog

test:
	go test ./...

install:
	go install ./cmd/clog

clean:
	rm -f clog
