build-all:
	cd cart && make build
	cd loms && make build
	cd notifier && make build

.PHONY: run-storage
run-storage:
	docker-compose up -d --wait loms-postgres-master loms-postgres-replica kafka0 kafka-init-topics kafka-ui
	cd loms && make migrate-postgres && make .setup-replication

run-monitoring:
	docker-compose up -d --wait jaeger prometheus grafana

run-all: build-all run-storage run-monitoring
	docker-compose build -q
	docker-compose up -d --force-recreate cart loms notifier


stop-all:
	docker-compose down
