package bootstrap

import (
	"context"
	"fmt"
	"log"
	postgres "stt/database/postgres"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresConnectionPool(env *Env) *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbHost := env.DBHost
	dbPort := env.DBPort
	dbUser := env.DBUser
	dbPass := env.DBPass
	dbName := env.DBName

	postgresDbURI := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	pool, err := postgres.InitConnectionPool(postgresDbURI)
	if err != nil {
		log.Fatal(err)
	}
	err = pool.Ping(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return pool
}

func ClosePostgresDbConnectionPool(pool *pgxpool.Pool) {
	pool.Close()
}
