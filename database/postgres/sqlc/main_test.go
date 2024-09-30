package db

import (
	"os"
	"stt/bootstrap"
	"testing"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	postgresConnectionPool := bootstrap.App().PostgresConnectionPool
	testQueries = New(postgresConnectionPool)
	os.Exit(m.Run())
}
