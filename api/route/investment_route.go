package route

import (
	"stt/api/controller"
	sv_interface "stt/services/interfaces"

	"github.com/gin-gonic/gin"
)

func InitInvestmentRouter(group *gin.RouterGroup, investmentService sv_interface.IInvestmentService) {
	ic := controller.InitInvestmentController(investmentService)

	group.GET("/investments", ic.Search)
	group.GET("/investment-overview/:id", ic.GetOverviewById)
	group.POST("/investment", ic.Create)
	group.PUT("/market-price", ic.UpdateMarketPrice)
}
