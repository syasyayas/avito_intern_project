
up:
	docker compose up --build -d

up-attached:
	docker compose up --build

down:
	docker compose down

clear_pg:
	docker rm avito_intern_project-postgres-1 && docker volume rm avito_intern_project_pg-data

.PHONY:
	: up down clear_pg
