package route

import (
	"stt/api/controller"
	sv_interface "stt/services/interfaces"

	"github.com/gin-gonic/gin"
)

func InitUserRouter(routerGroup *gin.RouterGroup, userService sv_interface.IUserService) {
	userController := controller.InitUserController(userService)
	routerGroup.POST("/users", userController.CreateUser)
}
