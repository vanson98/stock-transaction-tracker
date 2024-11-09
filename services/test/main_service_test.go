package service_test

import (
	"os"
	"stt/bootstrap"
	db "stt/database/postgres/sqlc"
	"stt/services"
	sv_interface "stt/services/interfaces"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var pgConnPool *pgxpool.Pool
var store db.IStore
var accService sv_interface.IAccountService
var userService sv_interface.IUserService
var tranService sv_interface.ITransactionService
var investmentService sv_interface.IInvestmentService

func TestMain(m *testing.M) {
	var timeout time.Duration = time.Second * 3
	env := bootstrap.NewEnv("../..")
	pgConnPool = bootstrap.NewPostgresConnectionPool(env)
	defer bootstrap.ClosePostgresDbConnectionPool(pgConnPool)

	store = db.NewStore(pgConnPool)
	accService = services.InitAccountService(store, timeout)
	userService = services.InitUserService(store)
	tranService = services.InitTransactionService(store)
	investmentService = services.InitInvestmentService(store, timeout)
	os.Exit(m.Run())
}
