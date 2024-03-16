.PHONY: build, coverage, lint
build:
	go build -o ./bin/app ./cmd

coverage:
	@go test -tags=unit -coverprofile=coverage.out  ./...
	@go tool cover -func coverage.out | grep total:
	@go tool cover -html coverage.out -o coverage.html
	@rm coverage.out

lint:
	@golangci-lint run

.PHONY: run-e2e
run-e2e:
	docker-compose build -q
	docker-compose up --force-recreate --detach
	-go test -count=1 -tags=e2e ./...
	docker-compose down