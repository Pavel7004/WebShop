all: build

name := shop

build:
	@echo "------------------"
	@echo "Building app...   "
	@echo "------------------"
	go build -o $(name) ./cmd/shop/shop.go

swag:
	@echo "------------------"
	@echo "Running swag...   "
	@echo "------------------"
	swag init --md ./ --pd -g server.go -d ./pkg/adapters/http

lint:
	@echo "------------------"
	@echo "Running linter... "
	@echo "------------------"
	golangci-lint run ./...

clear:
	rm $(name) *.out

clean:
	go clean -testcache
	go clean -cache

.PHONY: all build swag clear clean jaeger lint
