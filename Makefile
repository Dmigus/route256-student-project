build-all:
	cd cart && make build
	cd loms && make build

.PHONY: run-postgres
run-postgres:
	docker-compose up -d --wait loms-postgres-master
	docker-compose up -d --wait loms-postgres-replica
	cd loms && make migrate-postgres && make .setup-replication


run-all: build-all run-postgres
	docker-compose build -q
	docker-compose up -d --force-recreate cart loms

stop-all:
	docker-compose down
