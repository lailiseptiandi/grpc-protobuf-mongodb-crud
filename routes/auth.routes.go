package routes

import (
	"context"
	"grcp-api-client-mongo/controllers"
	"grcp-api-client-mongo/middleware"
	"grcp-api-client-mongo/services"
	"text/template"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthController struct {
	mongoDB *mongo.Database
}

func NewAuthControllerRoute(mongoDB *mongo.Database) AuthController {
	return AuthController{mongoDB}
}

func (r *AuthController) AuthRoute(rg *gin.RouterGroup) {

	ctx := context.TODO()
	var tmp *template.Template

	authService := services.NewAuthService(r.mongoDB.Collection("users"), ctx)
	userService := services.NewUserService(r.mongoDB.Collection("users"), ctx)
	authController := controllers.NewAuthController(authService, userService, ctx, r.mongoDB.Collection("users"), tmp)
	router := rg
	router.POST("/login", authController.LoginUser)
	router.POST("/register", authController.RegiserUser)
	router.GET("/refresh_token", authController.RefreshAccessToken)
	router.GET("/logout", middleware.AuthMiddleware(userService), authController.LogoutUser)
	router.POST("/forgot-password", authController.ForgotPassword)
	router.POST("/reset-password", authController.ResetPassword)

}
