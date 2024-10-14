package controller

import (
	"net/http"
	apimodels "stt/api/models"
	db "stt/database/postgres/sqlc"
	sv_interface "stt/services/interfaces"

	"github.com/gin-gonic/gin"
)

type userController struct {
	userService sv_interface.IUserService
}

func InitUserController(us sv_interface.IUserService) userController {
	return userController{
		userService: us,
	}
}

func (uc userController) CreateUser(ctx *gin.Context) {
	requestModel := apimodels.CreateUserRequest{}
	err := ctx.ShouldBindBodyWithJSON(&requestModel)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := uc.userService.CreateNew(ctx.Request.Context(), db.CreateUserParams{
		Username:       requestModel.Username,
		HashedPassword: requestModel.Password,
		FullName:       requestModel.FullName,
		Email:          requestModel.Email,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, user)
}
