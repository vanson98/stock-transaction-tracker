runapp:
	go run ./cmd/main.go
postgres:
	docker run --name postgres_db -e POSTGRES_USER=root -e POSTGRES_PASSWORD=123456 -p 5432:5432 -v pgdata:/var/lib/postgresql/data -d postgres
createdb:
	docker exec -it  postgres_db createdb --username=root --owner=root stock_tracker
dropdb:
	docker exec -it  postgres_db dropdb stock_tracker
migrateup:
	migrate -path database/postgres/migration --database "postgresql://root:123456@localhost:5432/stock_tracker?sslmode=disable" --verbose up
migratedown:
	migrate -path database/postgres/migration --database "postgresql://root:123456@localhost:5432/stock_tracker?sslmode=disable" --verbose down
sqlc: 
	sqlc generate
.PHONY: sqlc runapp postgres createdb dropdb migrateup migratedown