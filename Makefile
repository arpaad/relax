.PHONY: test lint tidy vet

VERSION := "v1.0.0"

test:
	go test ./... -race -v

lint:
	golangci-lint run ./...

tidy:
	go mod tidy

vet:
	go vet ./...

refresh-proxy:
	@echo "Click:"
	@echo "https://proxy.golang.org/github.com/luckyman42/relax/@v/${VERSION}.info"
	@echo "Then click:"
	@echo "https://pkg.go.dev/github.com/luckyman42/relax@${VERSION}"