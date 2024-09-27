package route

import (
	"stt/bootstrap"
	db "stt/database/postgres/sqlc"
	"time"

	"github.com/gin-gonic/gin"
)

func Setup(env *bootstrap.Env, timeout time.Duration, q *db.Queries, gin *gin.Engine) {
	//publicRouter := gin.Group("")

	// All public APIs
	protectedRouter := gin.Group("")
	NewInvestmentRouter(env, timeout, q, protectedRouter)
}
