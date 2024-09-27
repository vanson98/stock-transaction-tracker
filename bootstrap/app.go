package bootstrap

import "github.com/jackc/pgx/v5/pgxpool"

type Application struct {
	Env                    *Env
	PostgresConnectionPool *pgxpool.Pool
}

func App() Application {
	app := &Application{}
	app.Env = NewEnv()
	app.PostgresConnectionPool = NewPostgresConnectionPool(app.Env)
	return *app
}

func (app *Application) CloseDbConnection() {
	ClosePostgresDbConnectionPool(app.PostgresConnectionPool)
}
