package routes

import (
	"context"
	"grcp-api-client-mongo/controllers"
	"grcp-api-client-mongo/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostRouteController struct {
	mongoDB *mongo.Database
}

func NewPostControllerRoute(mongoDB *mongo.Database) PostRouteController {
	return PostRouteController{mongoDB}
}

func (r *PostRouteController) PostRoute(rg *gin.RouterGroup) {

	ctx := context.TODO()
	postService := services.NewPostService(r.mongoDB.Collection("posts"), ctx)
	postController := controllers.NewPostController(postService)
	router := rg.Group("/post")

	router.GET("/", postController.ListPost)
	router.GET("/:postId", postController.FindPostById)
	router.POST("/", postController.CreatePost)
	router.PATCH("/:postId", postController.UpdatePost)
	router.DELETE("/:postId", postController.DeletePost)
}
