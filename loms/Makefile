.PHONY: build
build:
	go build -o ./bin/app ./cmd

.PHONY: coverage
coverage:
	@go test -tags=unit,integration -coverprofile=coverage.out -parallel=10 ./...
	@go tool cover -html coverage.out -o coverage.html
	@rm coverage.out

race-detect:
	@go test -tags=unit,integration -count=1 -parallel=10 -race ./...

.PHONY:tutor-lint
tutor-lint:
	golangci-lint run --config=.golangci_tutor.yaml --new-from-rev=origin/master

.PHONY:tutor-lint-fix
tutor-lint-fix:
	golangci-lint run --config=.golangci_tutor.yaml --new-from-rev=origin/master --fix

.PHONY: lint
lint:
	@golangci-lint run

.PHONY: sqlc
sqlc:
	sqlc generate --file=./sqlc/modifier.yaml
	sqlc generate --file=./sqlc/reader.yaml

POSTGRES_PASSWORD:=$(shell cat ./db_password.txt)
.PHONY: migrate-postgres
migrate-postgres:
	@goose -dir ./migrations postgres "host=127.0.0.1 port=5432 user=postgres password=$(POSTGRES_PASSWORD) dbname=postgres sslmode=disable" up
	@goose -dir ./migrations postgres "host=127.0.0.1 port=5433 user=postgres password=$(POSTGRES_PASSWORD) dbname=postgres sslmode=disable" up

.PHONY: .setup-replication
.setup-replication:
	@docker exec $(shell docker ps -lq -f name=loms-postgres-master) psql -U postgres -c "CREATE PUBLICATION alltables FOR ALL TABLES;"  > /dev/null
	@docker exec $(shell docker ps -lq -f name=loms-postgres-replica) psql -U postgres -c "CREATE SUBSCRIPTION mysub CONNECTION 'host=loms-postgres-master port=5432 user=postgres password=$(POSTGRES_PASSWORD) dbname=postgres' PUBLICATION alltables WITH (synchronous_commit='on');" > /dev/null
	@docker exec $(shell docker ps -lq -f name=loms-postgres-replica) psql -U postgres -c "ALTER DATABASE postgres SET default_transaction_read_only to ON;"  > /dev/null


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

# Устанавливаем proto описания validate
vendor-proto/validate:
	git clone -b main --single-branch --depth=2 --filter=tree:0 \
		https://github.com/bufbuild/protoc-gen-validate vendor-proto/tmp && \
		cd vendor-proto/tmp && \
		git sparse-checkout set --no-cone validate &&\
		git checkout
		mkdir -p vendor-proto/validate
		mv vendor-proto/tmp/validate vendor-proto/
		rm -rf vendor-proto/tmp

# Устанавливаем proto описания google/googleapis
vendor-proto/google/api:
	git clone -b master --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/googleapis/googleapis vendor-proto/googleapis && \
 	cd vendor-proto/googleapis && \
	git sparse-checkout set --no-cone google/api && \
	git checkout
	mkdir -p  vendor-proto/google
	mv vendor-proto/googleapis/google/api vendor-proto/google
	rm -rf vendor-proto/googleapis

# Устанавливаем proto описания protoc-gen-openapiv2/options
vendor-proto/protoc-gen-openapiv2/options:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/grpc-ecosystem/grpc-gateway vendor-proto/grpc-ecosystem && \
 	cd vendor-proto/grpc-ecosystem && \
	git sparse-checkout set --no-cone protoc-gen-openapiv2/options && \
	git checkout
	mkdir -p vendor-proto/protoc-gen-openapiv2
	mv vendor-proto/grpc-ecosystem/protoc-gen-openapiv2/options vendor-proto/protoc-gen-openapiv2
	rm -rf vendor-proto/grpc-ecosystem


# Вендоринг внешних proto файлов
.vendor-proto: vendor-proto/google/protobuf vendor-proto/validate vendor-proto/google/api vendor-proto/protoc-gen-openapiv2/options

bin/protoc-gen-go:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1

bin/protoc-gen-go-grpc:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0

bin/protoc-gen-validate:
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@v1.0.4

bin/protoc-gen-grpc-gateway:
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.19.1

bin/protoc-gen-openapiv2:
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.19.1


.bin-deps: bin/protoc-gen-go bin/protoc-gen-go-grpc bin/protoc-gen-validate bin/protoc-gen-grpc-gateway bin/protoc-gen-openapiv2


LOMS_PROTO_PATH := api/loms
OUT_DIR := internal/controllers/grpc/protoc
SWAGGER_OUT_DIR := api/openapiv2

.PHONY: protoc-generate-service
protoc-generate-service: .bin-deps .vendor-proto
	mkdir -p $(OUT_DIR)
	mkdir -p $(SWAGGER_OUT_DIR)
	protoc \
	-I $(LOMS_PROTO_PATH) \
	-I vendor-proto \
	--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go \
	--go_out=$(OUT_DIR) \
	--go_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc \
	--go-grpc_out=$(OUT_DIR)  \
	--go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-validate=$(LOCAL_BIN)/protoc-gen-validate \
    --validate_out="lang=go,paths=source_relative:$(OUT_DIR)" \
    --plugin=protoc-gen-grpc-gateway=$(LOCAL_BIN)/protoc-gen-grpc-gateway \
	--grpc-gateway_out $(OUT_DIR) \
	--grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true \
	--plugin=protoc-gen-openapiv2=$(LOCAL_BIN)/protoc-gen-openapiv2 \
    --openapiv2_out $(SWAGGER_OUT_DIR) \
    --openapiv2_opt logtostderr=true \
	$(LOMS_PROTO_PATH)/v*/*.proto
	go mod tidy

