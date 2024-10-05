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
	pgConnPool = bootstrap.App("../../..").PostgresConnectionPool
	testQueries = New(pgConnPool)
	os.Exit(m.Run())
}
