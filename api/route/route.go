package route

import (
	"stt/api/middleware"
	"stt/bootstrap"
	db "stt/database/postgres/sqlc"
	"stt/services"
	"time"

	"github.com/gin-gonic/gin"
)

func Setup(env *bootstrap.Env, timeout time.Duration, store db.IStore, engine *gin.Engine) {
	middleware.UseCors(engine)
	//publicRouter := gin.Group("")

	accountService := services.InitAccountService(store, timeout)
	investmentService := services.InitInvestmentService(store, timeout)
	userService := services.InitUserService(store)

	// All protected APIs
	protectedRouter := engine.Group("")
	InitInvestmentRouter(protectedRouter, &investmentService)
	InitAccountRouter(protectedRouter, accountService)
	InitUserRouter(protectedRouter, userService)
}
