build-all:
	cd cart && make build
	cd loms && make build

.PHONY: run-postgres
run-postgres:
	docker-compose up -d --wait loms-postgres-master
	docker-compose up -d --wait loms-postgres-replica
	cd loms && make migrate-postgres && make .setup-replication

run-kafka:
	docker-compose up -d --wait kafka0 kafka-init-topics kafka-ui

run-all: build-all run-postgres run-kafka
	docker-compose build -q
	docker-compose up -d --force-recreate cart loms notifier


stop-all:
	docker-compose down
