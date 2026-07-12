.PHONY: build test lint clean hooks debt-scan generate

build:
	go build ./...

generate:
	go run ./cmd/asdt-schemagen

test:
	go test ./...

lint:
	golangci-lint run ./...

clean:
	go clean ./...

hooks:
	go install github.com/evilmartians/lefthook@latest
	lefthook install

debt-scan:
	grep -rn 'asdt:ceiling' --include='*.go' --include='*.md' . || true
