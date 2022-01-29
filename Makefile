BIN := "./bin/banner"
DOCKER_IMG="banner:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/banner

run: build
	$(BIN) -config ./configs/banner.json

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

generate:
	rm -rf api/pb
	mkdir -p api/pb

	protoc \
	--proto_path=api/ \
	--go_out=api/pb \
	--go-grpc_out=api/pb \
	api/*.proto

	protoc -I . --grpc-gateway_out api/pb\
    --grpc-gateway_opt logtostderr=true \
    --grpc-gateway_opt generate_unbound_methods=true \
    --proto_path=api/ \
    api/BannerService.proto

	protoc -I . --grpc-gateway_out api/pb\
    --grpc-gateway_opt logtostderr=true \
    --grpc-gateway_opt generate_unbound_methods=true \
    --proto_path=api/ \
    api/BannerService.proto

	protoc -I . \
    --go_out=":api/pb" \
    --validate_out="lang=go:api/pb" \
     --proto_path=api/ \
    api/BannerService.proto