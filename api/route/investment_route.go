package route

import (
	"stt/api/controller"
	"stt/domain"

	"github.com/gin-gonic/gin"
)

func InitInvestmentRouter(group *gin.RouterGroup, investmentService *domain.IInvestmentService) {
	ic := controller.InvestmentController{
		InvestmentService: *investmentService,
	}

	group.GET("/investments", ic.GetAll)
	group.POST("/investment", func(c *gin.Context) {})
}
