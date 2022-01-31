BIN := "./bin/rotator"
DOCKER_IMG="rotator:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/rotator

run: build
	$(BIN) -config ./configs/rotator.json

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
	go test -v -race -count=100 -timeout=1m ./...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

migration:
	goose -dir=migrations postgres "user=alex_molokov password=alex_molokov dbname=banner sslmode=disable" up

migration-reverse:
	goose -dir=migrations postgres "user=alex_molokov password=alex_molokov dbname=banner sslmode=disable" down

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
    api/RotatorService.proto

	protoc -I . --grpc-gateway_out api/pb\
    --grpc-gateway_opt logtostderr=true \
    --grpc-gateway_opt generate_unbound_methods=true \
    --proto_path=api/ \
    api/RotatorService.proto

	protoc -I . \
    --go_out=":api/pb" \
    --validate_out="lang=go:api/pb" \
     --proto_path=api/ \
    api/RotatorService.proto