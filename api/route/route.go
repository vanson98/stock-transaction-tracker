package route

import (
	"stt/bootstrap"
	db "stt/database/postgres/sqlc"
	"time"

	"github.com/gin-gonic/gin"
)

func Setup(env *bootstrap.Env, timeout time.Duration, store db.IStore, gin *gin.Engine) {
	//publicRouter := gin.Group("")

	// All protected APIs
	protectedRouter := gin.Group("")

	NewInvestmentRouter(env, timeout, store, protectedRouter)
	NewAccountRouter(store, timeout, protectedRouter)
}
