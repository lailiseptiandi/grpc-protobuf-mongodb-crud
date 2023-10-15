package controllers

import (
	"grcp-api-client-mongo/models"
	"grcp-api-client-mongo/services"
	"grcp-api-client-mongo/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) UserController {
	return UserController{userService}
}

func (uc *UserController) GetMeProfile(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(*models.DBResponseUser)

	resp := utils.ResponseSuccess(currentUser, "Successfully get profile")
	ctx.JSON(http.StatusOK, resp)
	return
}
