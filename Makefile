# Minimal System Metrics Agent (msma) Makefile
# ==================================================

.PHONY: all fmt vet test build run clean

all: fmt vet test

fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	go test -v -race -count=1 ./...

build:
	go build -o msma .

run:
	go run .

clean:
	rm -f msma
