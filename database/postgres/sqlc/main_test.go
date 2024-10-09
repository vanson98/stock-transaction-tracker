package db

import (
	"os"
	"stt/bootstrap"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testQueries *Queries
var pgConnPool *pgxpool.Pool

func TestMain(m *testing.M) {
	env := bootstrap.NewEnv("../../..")
	pgConnPool = bootstrap.NewPostgresConnectionPool(env)
	testQueries = New(pgConnPool)
	os.Exit(m.Run())
}
