package services

import (
	"os"
	"stt/bootstrap"
	db "stt/database/postgres/sqlc"
	"stt/domain"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var pgConnPool *pgxpool.Pool
var store db.IStore
var accService domain.IAccountService

func TestMain(m *testing.M) {
	var timeout time.Duration = 3
	pgConnPool = bootstrap.App("..").PostgresConnectionPool
	store = db.NewStore(pgConnPool)

	accService = InitAccountService(store, timeout)
	os.Exit(m.Run())
}
