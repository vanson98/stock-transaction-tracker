package route

import (
	"stt/api/controller"
	sv_interface "stt/services/interfaces"

	"github.com/gin-gonic/gin"
)

func InitAccountRouter(routerGroup *gin.RouterGroup,
	accountService sv_interface.IAccountService,
	userService sv_interface.IUserService) {
	accountController := controller.AccountController{
		AccountService: accountService,
		UserService:    userService,
	}

	routerGroup.POST("/accounts", accountController.CreateNewAccount)
	routerGroup.GET("/account-overview", accountController.GetAccountOverview)
	routerGroup.PUT("/account-transfer", accountController.TransferMoney)
	routerGroup.GET("/accounts", accountController.GetListAccount)
	routerGroup.GET("/account-search", accountController.GetAllAccountOverview)
}
