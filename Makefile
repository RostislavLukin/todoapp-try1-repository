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