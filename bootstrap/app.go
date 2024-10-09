package bootstrap

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct {
	Env                    *Env
	PostgresConnectionPool *pgxpool.Pool
	Engine                 *gin.Engine
}

func NewServerApp(envPath string) Application {
	app := Application{}
	app.Env = NewEnv(envPath)
	app.PostgresConnectionPool = NewPostgresConnectionPool(app.Env)
	app.Engine = gin.Default()
	return app
}

func (app *Application) CloseDbConnection() {
	ClosePostgresDbConnectionPool(app.PostgresConnectionPool)
}
