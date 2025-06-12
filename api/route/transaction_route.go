package route

import (
	"stt/api/controller"
	sv_interface "stt/services/interfaces"

	"github.com/gin-gonic/gin"
)

func InitTransactionRouter(group *gin.RouterGroup, transactionService sv_interface.ITransactionService) {
	transactionController := controller.InitTransactionController(transactionService)
	group.GET("/transactions", transactionController.GetPaging)
	group.POST("/transaction", transactionController.CreateNewTransaction)
	group.POST("transaction/importing", transactionController.AddTransaction)
}
