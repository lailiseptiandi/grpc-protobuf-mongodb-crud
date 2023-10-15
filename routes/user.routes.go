package routes

import (
	"context"
	"grcp-api-client-mongo/controllers"
	"grcp-api-client-mongo/middleware"
	"grcp-api-client-mongo/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRouteController struct {
	mongoDB *mongo.Database
}

func NewRouteUserController(mongoDB *mongo.Database) UserRouteController {
	return UserRouteController{mongoDB}
}

func (uc *UserRouteController) UserRoute(rg *gin.RouterGroup) {
	ctx := context.TODO()
	userService := services.NewUserService(uc.mongoDB.Collection("users"), ctx)
	userController := controllers.NewUserController(userService)
	router := rg.Group("/users")
	router.Use(middleware.AuthMiddleware(userService))
	router.GET("/profile", userController.GetMeProfile)
}
