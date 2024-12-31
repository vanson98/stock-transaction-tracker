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
	routerGroup.GET("/asset-catalog", accountController.GetAssetCatalog)
	routerGroup.PUT("/account-transfer", accountController.TransferMoney)
	routerGroup.GET("/accounts", accountController.GetListAccount)
	routerGroup.GET("/account-search", accountController.GetAllAccountOverview)
}
