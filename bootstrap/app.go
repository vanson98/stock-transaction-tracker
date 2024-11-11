package bootstrap

import (
	controler_validator "stt/api/validator"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
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
	// custom param validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", controler_validator.ValidCurrency)
		v.RegisterValidation("trade", controler_validator.ValidTradeType)
	}
	return app
}

func (app *Application) CloseDbConnection() {
	ClosePostgresDbConnectionPool(app.PostgresConnectionPool)
}
