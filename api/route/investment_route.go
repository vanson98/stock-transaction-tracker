package route

import (
	sv_interface "stt/services/interfaces"

	"github.com/gin-gonic/gin"
)

func InitInvestmentRouter(group *gin.RouterGroup, investmentService *sv_interface.IInvestmentService) {
	// ic := controller.InvestmentController{
	// 	InvestmentService: *investmentService,
	// }

	//group.GET("/investments", ic.GetAll)
	//group.POST("/investment", func(c *gin.Context) {})
}
