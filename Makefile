all: lint

build:
	@echo "------------------"
	@echo "Building app...   "
	@echo "------------------"
	go build cmd/shop/shop.go

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

jaeger:
	docker run -dp 6831:6831/udp -p 16686:16686 jaegertracing/all-in-one:latest

clear:
	rm shop *.out

clean:
	go clean -testcache
	go clean -cache

.PHONY: all build swag clear clean jaeger lint
