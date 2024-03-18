# Используем bin в текущей директории для установки плагинов protoc
LOCAL_BIN:=$(CURDIR)/bin

#.PHONY: .vendor-rm
#.vendor-rm:
#	rm -rf vendor-proto

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

#.PHONY: .bin-deps
.bin-deps: bin/protoc-gen-go bin/protoc-gen-go-grpc


LOMS_PROTO_PATH:=api
OUT_DIR := "internal/protoc"

.PHONY: .protoc-generate
.protoc-generate: .bin-deps .vendor-proto
	mkdir -p $(OUT_DIR)
	protoc \
	-I . \
	-I vendor-proto \
	--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go \
	--go_out= \
	--go_opt=module="route256.ozon.ru/project" \
	--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc \
	--go-grpc_out=. \
	--go-grpc_opt=module="route256.ozon.ru/project" \
	$(LOMS_PROTO_PATH)/**/**/*.proto
	go mod tidy


