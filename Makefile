include .env
export

SERVICE?=$(SERVICE_WALLET)
SERVICE_DB?=$(WALLET_DB)
DB_URL?=postgres://$(PAYMENTS_USER):$(PAYMENTS_PASSWORD)@localhost:5432/$(SERVICE_DB)?sslmode=disable
MIGRATE := migrate -path db/$(SERVICE)/migrations -database "$(DB_URL)"

.PHONY: migrate-up migrate-down migrate-version migrate-new docker-run-sql-file docker-run-sql-query

migrate-up:
	@$(MIGRATE) up
migrate-down:
	@$(MIGRATE) down 1
migrate-version:
	@$(MIGRATE) version
migrate-new:
	@test -n "$(name)" || { echo "usage: make migrate-new name=add_name_of_sql_files_of_migration"; exit 1; }
	@migrate create -ext sql -dir db/$(SERVICE)/migrations -seq "$(name)"

docker-run-sql-file:
	@test -n "$(file)" || { echo "usage: make docker-run-sql-file name=scratch.sql database=database_name"; exit 1; }
	@test -n "$(database)" || { echo "usage: make docker-run-sql-file name=scratch.sql database=database_name"; exit 1; }
	@docker compose exec -T postgres psql -U payments -d $(database) -q < "$(file)"

docker-run-sql-query:
	@test -n "$(query)" || { echo "usage: make docker-run-sql-query query='SELECT 1;'"; exit 1; }
	@docker compose exec -T postgres psql -U payments -d payment_db -c "$(query)"