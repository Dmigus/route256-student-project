build-all:
	cd cart && make build
	cd loms && make build

run-all: build-all
	docker-compose build -q
	docker-compose up --force-recreate
