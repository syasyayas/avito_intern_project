
dockerrun:
	docker compose up --build
dockerdown:
	docker compose down

clear_pg:
	docker rm avito_intern_project-postgres-1 && docker volume rm avito_intern_project_pg-data

.PHONY:
	: dockerrun dockerdown clear_pg
