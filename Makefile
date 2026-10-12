include .env
export

export PROJECT_ROOT=$(shell pwd)

.PHONY: env-up env-down env-port-forward env-port-close env-cleanup 

env-up:
	mkdir -p out/pgdata
	docker compose up -d todoapp-postgres

env-down:
	docker compose down todoapp-postgres

env-port-forward:
	@docker compose up -d port-forward

env-port-close:
	@docker compose down port-forward

env-cleanup:
	@read -p "Очистка всех volume-файлов окружения? Опасность потери данных [y/N]: " ans; \
	if [ "$$ans" = "y" ] || [ "$$ans" = "Y" ]; then \
		docker compose down todoapp-postgres port-forward && \
		rm -rf ${PROJECT_ROOT}/out/pgdata && \
		mkdir -p out && \
		echo "Файлы окружения очищены!"; \
	else \
		echo "Очистка окружения отменена"; \
	fi

.PHONY: migrate-creat migrate-up migrate-down migrate-action

migrate-creat:
	@if [ -z "$(seq)" ]; then \
		echo "Отсутствует необходимый параметр seq. Пример: make migrate-creat seq=init"; \
		exit 1; \
	fi
	mkdir -p migrations
	docker compose run --rm todoapp-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down
	
migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Отсутствует необходимый параметр seq. Пример: make migrate-action action=upt"; \
		exit 1; \
	fi
	@docker compose run --rm todoapp-postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@todoapp-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		${action}


todoapp-run:
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
	export POSTGRES_HOST=localhost && \
	export POSTGRES_PORT=$${POSTGRES_FORWARD_PORT:-5433} && \
	go mod tidy && \
	go run ${PROJECT_ROOT}/cmd/todoapp/main.go
