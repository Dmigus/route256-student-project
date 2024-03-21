build-all:
	cd cart && GOOS=linux GOARCH=amd64 make build
	cd loms && GOOS=linux GOARCH=amd64 make build

run-all: build-all
	docker-compose build -q
	docker-compose up --force-recreate
