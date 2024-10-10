package route

import (
	"stt/bootstrap"
	db "stt/database/postgres/sqlc"
	"stt/services"
	"time"

	"github.com/gin-gonic/gin"
)

func Setup(env *bootstrap.Env, timeout time.Duration, store db.IStore, gin *gin.Engine) {
	//publicRouter := gin.Group("")

	accountService := services.InitAccountService(store, timeout)
	investmentService := services.InitInvestmentService(store, timeout)
	// All protected APIs
	protectedRouter := gin.Group("")

	InitInvestmentRouter(protectedRouter, &investmentService)
	InitAccountRouter(protectedRouter, accountService)

}
