server:
	go run ./cmd/main.go

postgres:
	docker run --name postgres_db -e POSTGRES_USER=root -e POSTGRES_PASSWORD=123456 -p 5432:5432 -v pgdata:/var/lib/postgresql/data -d postgres

createdb:
	docker exec -it  postgres_db createdb --username=root --owner=root stock_tracker

dropdb:
	docker exec -it  postgres_db dropdb stock_tracker

migrateup:
	migrate -path database/postgres/migration -database "postgresql://root:123456@localhost:5432/stock_tracker?sslmode=disable" -verbose up

migrateup1:
	migrate -path database/postgres/migration -database "postgresql://root:123456@localhost:5432/stock_tracker?sslmode=disable" -verbose up 1

migratedown:
	migrate -path database/postgres/migration -database "postgresql://root:123456@localhost:5432/stock_tracker?sslmode=disable" -verbose down

migratedown1:
	migrate -path database/postgres/migration -database "postgresql://root:123456@localhost:5432/stock_tracker?sslmode=disable" -verbose down 1

new_migration:
	migrate create -ext sql -dir database/postgres/migration -seq $(name)

sqlc: 
	sqlc generate

service_test: 
	go test  -coverpkg=./services -v -cover -coverprofile=coverage.out ./services/test

cov2lcov:
	gcov2lcov -infile=coverage.out -outfile=lcov.info 

controller_test: 
	go test -coverpkg=./api/controller -v -cover ./api/test

mock:
	mockgen -package mock_service  -destination services/mock/service_mock.go stt/services/interfaces IAccountService,IInvestmentService,IUserService
.PHONY: sqlc runapp postgres createdb dropdb migrateup migratedown server mock service_test controller_test cov2lcov