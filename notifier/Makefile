.PHONY: build
build:
	go build -o ./bin/app ./cmd

.PHONY:tutor-lint
tutor-lint:
	golangci-lint run --config=.golangci_tutor.yaml --new-from-rev=origin/master

.PHONY:tutor-lint-fix
tutor-lint-fix:
	golangci-lint run --config=.golangci_tutor.yaml --new-from-rev=origin/master --fix

