1. run : go install -tags 'potgres mongodb sqlite etc(becouse it overrides older migrate)' go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
2. migrate create -ext sql -dir db/wallet/migrations -seq init_ledger
3. migrate -path db/wallet/migrations \
  -database "postgres://payments:Payments%40123@localhost:5432/wallet_db?sslmode=disable" up
4. docker compose exec postgres  psql -U payments -d wallet_db -c 'SELECT * FROM schema_migrations' -- to Connect to db in postgres and talk to postgres we alwyas need to connect to some db