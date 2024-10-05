package bootstrap

import "github.com/jackc/pgx/v5/pgxpool"

type Application struct {
	Env                    *Env
	PostgresConnectionPool *pgxpool.Pool
}

func App(envPath string) Application {
	app := &Application{}
	app.Env = NewEnv(envPath)
	app.PostgresConnectionPool = NewPostgresConnectionPool(app.Env)
	return *app
}

func (app *Application) CloseDbConnection() {
	ClosePostgresDbConnectionPool(app.PostgresConnectionPool)
}
