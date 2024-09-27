package route

import (
	"stt/api/controller"
	"stt/bootstrap"
	db "stt/database/postgres/sqlc"
	"stt/repositories"
	"stt/services"
	"time"

	"github.com/gin-gonic/gin"
)

func NewInvestmentRouter(env *bootstrap.Env, timeout time.Duration, queries *db.Queries, group *gin.RouterGroup) {
	ir := repositories.InitInvestmentRepository(queries)
	ic := controller.InvestmentController{
		InvestmentService: services.InitInvestmentService(ir, timeout),
	}

	group.GET("/investments", ic.GetAll)
	group.POST("/investment", func(c *gin.Context) {})
}
