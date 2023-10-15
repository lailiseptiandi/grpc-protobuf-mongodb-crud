package controllers

import (
	"context"
	"fmt"
	"grcp-api-client-mongo/config"
	"grcp-api-client-mongo/models"
	"grcp-api-client-mongo/services"
	"grcp-api-client-mongo/utils"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthController struct {
	authService services.AuthService
	userService services.UserService
	ctx         context.Context
	collection  *mongo.Collection
	temp        *template.Template
}

func NewAuthController(authService services.AuthService, userService services.UserService, ctx context.Context, collection *mongo.Collection, temp *template.Template) AuthController {
	return AuthController{authService, userService, ctx, collection, temp}
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

func (ac *AuthController) ForgotPassword(ctx *gin.Context) {
	var userCredential *models.ForgotPasswordInput

	err := ctx.ShouldBindJSON(&userCredential)

	if err != nil {
		resp := utils.ResponseError(nil, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	message := "You will receive a reset email if user with that email exist"

	user, err := ac.userService.FindUserByEmail(userCredential.Email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			resp := utils.ResponseError(nil, message)
			ctx.JSON(http.StatusOK, resp)
			return
		}
		resp := utils.ResponseError(nil, err.Error())
		ctx.JSON(http.StatusBadGateway, resp)
		return
	}

	if !user.Verified {
		resp := utils.ResponseError(nil, "Account not verified")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	config, err := config.LoadConfig(".")

	if err != nil {
		log.Fatal("could not load config : ", err)
	}

	// generate code
	resetToken := randstr.String(20)

	passwordResetToken := utils.Encode(resetToken)

	// update user
	query := bson.D{{Key: "email", Value: strings.ToLower(userCredential.Email)}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "passwordResetToken", Value: passwordResetToken}, {Key: "passwordResetAt", Value: time.Now().Add(time.Minute * 15)}}}}
	result, err := ac.collection.UpdateOne(ac.ctx, query, update)

	if result.MatchedCount == 0 {
		resp := utils.ResponseError(nil, "There was an error sending email")
		ctx.JSON(http.StatusBadGateway, resp)
		return
	}

	if err != nil {
		resp := utils.ResponseError(nil, err.Error())
		ctx.JSON(http.StatusForbidden, resp)
		return
	}
	var firstName = user.Name

	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[1]
	}

	// ? Send Email
	emailData := utils.EmailData{
		URL:       config.Origin + "/resetpassword/" + resetToken,
		FirstName: firstName,
		Subject:   "Your password reset token (valid for 10min)",
	}

	err = utils.SendEmail(user, &emailData, ac.temp, "resetPassword.html")
	if err != nil {
		resp := utils.ResponseError(nil, "There was an error sending email")
		ctx.JSON(http.StatusBadGateway, resp)
		return
	}
	resp := utils.ResponseSuccess(nil, message)
	ctx.JSON(http.StatusOK, resp)
	return
}

func (ac *AuthController) ResetPassword(ctx *gin.Context) {
	resetToken := ctx.Params.ByName("resetToken")
	var userCredential *models.ResetPasswordInput

	if err := ctx.ShouldBindJSON(&userCredential); err != nil {
		resp := utils.ResponseError(nil, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	if userCredential.Password != userCredential.PasswordConfirm {
		resp := utils.ResponseError(nil, "Passwords do not match")
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	hashedPassword, _ := utils.HashPassword(userCredential.Password)

	passwordResetToken := utils.Encode(resetToken)

	// Update User in Database
	query := bson.D{{Key: "passwordResetToken", Value: passwordResetToken}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "password", Value: hashedPassword}}}, {Key: "$unset", Value: bson.D{{Key: "passwordResetToken", Value: ""}, {Key: "passwordResetAt", Value: ""}}}}
	result, err := ac.collection.UpdateOne(ac.ctx, query, update)

	if result.MatchedCount == 0 {
		resp := utils.ResponseError(nil, "Token is invalid or has expired")
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	if err != nil {
		resp := utils.ResponseError(nil, err.Error())
		ctx.JSON(http.StatusForbidden, resp)
		return
	}

	ctx.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "", -1, "/", "localhost", false, true)

	resp := utils.ResponseError(nil, "Password data updated successfully")
	ctx.JSON(http.StatusOK, resp)
	return
}
