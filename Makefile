.PHONY: build run fmt vet test check

build:
	go build -o bin/spent ./cmd/spent

run:
	go run ./cmd/spent

fmt:
	gofmt -w .

vet:
	go vet ./...

test:
	go test ./...

check:
	@test -z "$$(gofmt -l .)" || (echo "gofmt needed:"; gofmt -l .; exit 1)
	go vet ./...
	go test ./...
