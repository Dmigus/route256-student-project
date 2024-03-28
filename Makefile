build-all:
	cd cart && make build
	cd loms && make build

run-all: build-all
	docker-compose build -q
	docker-compose up -d --wait loms-postgres
	cd loms && make migrate-postgres
	docker-compose up -d --force-recreate cart loms

stop-all:
	docker-compose down
