build:
	go build -o ./bin/app ./cmd

coverage:
	@go test -coverprofile=coverage.out ./internal/usecases/... ./internal/providers/repository/...> /dev/null
	@go tool cover -func coverage.out | grep total:
	@rm coverage.out