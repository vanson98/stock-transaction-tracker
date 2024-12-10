package route

import (
	"stt/api/controller"
	sv_interface "stt/services/interfaces"

	"github.com/gin-gonic/gin"
)

func InitAccountRouter(routerGroup *gin.RouterGroup, accountService sv_interface.IAccountService) {
	accountController := controller.AccountController{
		AccountService: accountService,
	}

	routerGroup.POST("/accounts", accountController.CreateNewAccount)
	routerGroup.PUT("/account-transfer", accountController.TransferMoney)
	routerGroup.GET("/accounts", accountController.GetListAccount)
	routerGroup.GET("/accounts/:id", accountController.GetAccountById)
	routerGroup.GET("/account-info/:id", accountController.GetAccountInfoById)
	routerGroup.GET("/account-search", accountController.GetAccoutPaging)
}
