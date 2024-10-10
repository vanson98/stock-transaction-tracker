package route

import (
	"stt/api/controller"
	"stt/domain"

	"github.com/gin-gonic/gin"
)

func InitAccountRouter(routerGroup *gin.RouterGroup, accountService domain.IAccountService) {
	accountController := controller.AccountController{
		AccountService: accountService,
	}

	routerGroup.POST("/accounts", accountController.CreateNewAccount)
	routerGroup.GET("/accounts/:id", accountController.GetAccountById)
	routerGroup.PUT("/account-transfer", accountController.TransferMoney)
}
