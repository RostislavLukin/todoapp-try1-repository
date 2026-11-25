include .env

export

export PROJECT_ROOT=$(shell pwd)

env-up:
	@docker compose up -d env-postgres

env-down:
	@docker compose down env-postgres

port-forwarder:
	@docker compose up -d port-forward

port-close:
	@docker compose down port-forward

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Отсутствует необходимый параметр 'seq'. Например: make migrate-create seq=init."; \
		exit 1; \
	fi; \
	docker compose run --rm postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-up:
	@make migrate-actioin action=up
migrate-down:
	@make migrate-action action=down
migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Отсутствует обязательный параметр 'action'. Пример: make migrate-action action=2"; \
		exit 1; \
	fi;
	@docker compose run --rm postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@env-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		"$(action)"