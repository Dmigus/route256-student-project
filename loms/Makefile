.PHONY: build
build:
	go build -o ./bin/app ./cmd

.PHONY: coverage
coverage:
	@go test -tags=unit -coverprofile=coverage.out  ./...
	@go tool cover -func coverage.out | grep total:
	@go tool cover -html coverage.out -o coverage.html
	@rm coverage.out

lint:
	@golangci-lint run

# Используем bin в текущей директории для установки плагинов protoc
LOCAL_BIN:=$(CURDIR)/bin

# Устанавливаем proto описания google/protobuf
vendor-proto/google/protobuf:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/protocolbuffers/protobuf vendor-proto/protobuf &&\
	cd vendor-proto/protobuf &&\
	git sparse-checkout set --no-cone src/google/protobuf &&\
	git checkout
	mkdir -p vendor-proto/google
	mv vendor-proto/protobuf/src/google/protobuf vendor-proto/google
	rm -rf vendor-proto/protobuf


# Вендоринг внешних proto файлов
.vendor-proto: vendor-proto/google/protobuf

bin/protoc-gen-go:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1

bin/protoc-gen-go-grpc:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0

.bin-deps: bin/protoc-gen-go bin/protoc-gen-go-grpc


LOMS_PROTO_PATH := api/loms
OUT_DIR := internal/controllers/protoc

.PHONY: .protoc-generate
.protoc-generate: .bin-deps .vendor-proto
	mkdir -p $(OUT_DIR)
	protoc \
	-I $(LOMS_PROTO_PATH) \
	-I vendor-proto \
	--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go \
	--go_out=$(OUT_DIR) \
	--go_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc \
	--go-grpc_out=$(OUT_DIR)  \
	--go-grpc_opt=paths=source_relative \
	$(LOMS_PROTO_PATH)/v*/*.proto
	go mod tidy


