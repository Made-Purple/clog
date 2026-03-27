.PHONY: build test install clean

build:
	go build -o clog ./cmd/clog

test:
	go test ./...

install:
	go install ./cmd/clog

clean:
	rm -f clog
