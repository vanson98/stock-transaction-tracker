package route

import (
	"stt/api/controller"
	db "stt/database/postgres/sqlc"
	"stt/services"
	"time"

	"github.com/gin-gonic/gin"
)

func NewAccountRouter(dbStore db.IStore, timeout time.Duration, routerGroup *gin.RouterGroup) {
	accountController := controller.AccountController{
		AccountService: services.InitAccountService(dbStore, timeout),
	}

	routerGroup.POST("/account", accountController.CreateNewAccount)
}
