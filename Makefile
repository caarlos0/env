setup:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh

verify:
	go test -v -cover ./...
	./bin/golangci-lint run --enable-all --disable=lll ./...

.DEFAULT_GOAL := verify
