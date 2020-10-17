.PHONY: build test run

build:
	go vet ./...
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -ldflags "-s -w" -o ./build/tracker ./cmd/tracker

test:
	go test -tags=unit -race -short `go list -tags=unit ./...`
