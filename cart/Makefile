.PHONY: build, coverage, lint
build:
	go build -o ./bin/app ./cmd

coverage:
	@go test -coverprofile=coverage.out ./internal/usecases/... ./internal/providers/repository/... ./tests/...> /dev/null
	@go tool cover -func coverage.out | grep total:
	@rm coverage.out

lint:
	@golangci-lint run