package controllers

import (
	"fmt"
	"grcp-api-client-mongo/config"
	"grcp-api-client-mongo/models"
	"grcp-api-client-mongo/services"
	"grcp-api-client-mongo/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthController struct {
	authService services.AuthService
	userService services.UserService
}

func NewAuthController(authService services.AuthService, userService services.UserService) AuthController {
	return AuthController{authService, userService}
}

func (ac *AuthController) RegiserUser(ctx *gin.Context) {
	var userRegister *models.RegiserUser

	err := ctx.ShouldBindJSON(&userRegister)
	if err != nil {
		resp := utils.ResponseError(nil, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	if userRegister.Password != userRegister.PasswordConfirm {
		resp := utils.ResponseError(nil, "Passwords do not match")
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	newUser, err := ac.authService.RegisterUser(userRegister)
	if err != nil {
		if strings.Contains(err.Error(), "email already exist") {
			resp := utils.ResponseError(nil, err.Error())
			ctx.JSON(http.StatusConflict, resp)
			return
		}

		resp := utils.ResponseError(nil, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	config, _ := config.LoadConfig(".")

	accessToken, err := utils.CreateToken(config.AccessTokenExpiresIn, newUser.ID, config.AccessTokenPrivateKey)

	if err != nil {
		resp := utils.ResponseError(nil, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	refreshToken, err := utils.CreateToken(config.RefreshTokenExpiresIn, newUser.ID, config.RefreshTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	ctx.SetCookie("access_token", accessToken, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refreshToken, config.RefreshTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	formatter := models.FormatterLoginRegister(newUser, accessToken, refreshToken)
	resp := utils.ResponseSuccess(formatter, "Successfully register user")
	ctx.JSON(http.StatusCreated, resp)
	return
}

func (ac *AuthController) LoginUser(ctx *gin.Context) {

	var userLogin *models.LoginUser
	err := ctx.ShouldBindJSON(&userLogin)
	if err != nil {
		resp := utils.ResponseError(nil, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	user, err := ac.userService.FindUserByEmail(userLogin.Email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			resp := utils.ResponseError(nil, "Invalid email or password")
			ctx.JSON(http.StatusBadRequest, resp)
			return
		}
		resp := utils.ResponseError(nil, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	if err := utils.VerifyPassword(user.Password, userLogin.Password); err != nil {
		resp := utils.ResponseError(nil, "Invalid email or password")
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	config, _ := config.LoadConfig(".")

	accessToken, err := utils.CreateToken(config.AccessTokenExpiresIn, user.ID, config.AccessTokenPrivateKey)

	if err != nil {
		resp := utils.ResponseError(nil, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	refreshToken, err := utils.CreateToken(config.RefreshTokenExpiresIn, user.ID, config.RefreshTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	ctx.SetCookie("access_token", accessToken, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refreshToken, config.RefreshTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, false)
	formatter := models.FormatterLoginRegister(user, accessToken, refreshToken)
	resp := utils.ResponseSuccess(formatter, "Successfully login")
	ctx.JSON(http.StatusOK, resp)
	return
}

func (ac *AuthController) RefreshAccessToken(ctx *gin.Context) {
	message := "could not refresh access token"

	cookie, err := ctx.Cookie("refresh_token")
	if err != nil {
		resp := utils.ResponseError(nil, message)
		ctx.AbortWithStatusJSON(http.StatusForbidden, resp)
		return
	}

	config, _ := config.LoadConfig(".")
	sub, err := utils.ValidateToken(cookie, config.RefreshTokenPublicKey)
	if err != nil {
		resp := utils.ResponseError(nil, err.Error())
		ctx.AbortWithStatusJSON(http.StatusForbidden, resp)
		return
	}

	user, err := ac.userService.FindUserById(fmt.Sprint(sub))
	if err != nil {
		resp := utils.ResponseError(nil, err.Error())
		ctx.AbortWithStatusJSON(http.StatusForbidden, resp)
		return
	}

	accessToken, err := utils.CreateToken(config.AccessTokenExpiresIn, user, config.AccessTokenPrivateKey)
	if err != nil {
		resp := utils.ResponseError(nil, err.Error())
		ctx.AbortWithStatusJSON(http.StatusForbidden, resp)
		return
	}
	ctx.SetCookie("access_token", accessToken, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	formatter := map[string]interface{}{
		"access_token": accessToken,
	}
	resp := utils.ResponseSuccess(formatter, "Successfully refresh token")
	ctx.JSON(http.StatusOK, resp)
	return

}

func (ac *AuthController) LogoutUser(ctx *gin.Context) {
	ctx.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "", -1, "/", "localhost", false, true)

	resp := utils.ResponseSuccess(nil, "Successfully Logout")
	ctx.JSON(http.StatusOK, resp)
	return
}
