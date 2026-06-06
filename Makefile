.PHONY: test lint tidy vet

VERSION := "v0.1.0"

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
	@echo "https://proxy.golang.org/github.com/arpaad/relax/@v/${VERSION}.info"
	@echo "Then click:"
	@echo "https://pkg.go.dev/github.com/arpaad/relax@${VERSION}"