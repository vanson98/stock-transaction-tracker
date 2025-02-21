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

	//publicRouter := gin.Group("")

	accountService := services.InitAccountService(store, timeout)
	investmentService := services.InitInvestmentService(store, timeout)
	userService := services.InitUserService(store)
	transactionService := services.InitTransactionService(store)

	// All protected APIs
	protectedRouter := engine.Group("")

	// set up middlewares
	middleware.UseCors(protectedRouter)
	middleware.UseTokenVerification(protectedRouter)

	// set up router matching pattern
	InitInvestmentRouter(protectedRouter, investmentService)
	InitAccountRouter(protectedRouter, accountService, userService)
	InitUserRouter(protectedRouter, userService)
	InitTransactionRouter(protectedRouter, transactionService)
}
