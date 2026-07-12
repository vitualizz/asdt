.PHONY: build test lint clean hooks debt-scan

build:
	go build ./...

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
