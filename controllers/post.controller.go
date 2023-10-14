package controllers

import (
	"grcp-api-client-mongo/models"
	"grcp-api-client-mongo/services"
	"grcp-api-client-mongo/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type PostController struct {
	postService services.PostService
}

func NewPostController(postService services.PostService) PostController {
	return PostController{postService}
}

func (pc *PostController) CreatePost(ctx *gin.Context) {
	var post *models.CreatePostRequest
	err := ctx.ShouldBindJSON(&post)
	if err != nil {
		resp := utils.ResponseError(nil, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	newPost, err := pc.postService.CreatePost(post)
	if err != nil {
		if strings.Contains(err.Error(), "title already exists") {
			resp := utils.ResponseError(nil, err.Error())
			ctx.JSON(http.StatusConflict, resp)
			return
		}
		resp := utils.ResponseError(nil, err.Error())
		ctx.JSON(http.StatusBadGateway, resp)
		return
	}
	resp := utils.ResponseSuccess(newPost, "Successfully created post")
	ctx.JSON(http.StatusCreated, resp)
	return
}

func (pc *PostController) UpdatePost(ctx *gin.Context) {
	var post *models.UpdatePost
	postId := ctx.Param("postId")

	err := ctx.ShouldBindJSON(&post)
	if err != nil {
		resp := utils.ResponseError(nil, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	updatePost, err := pc.postService.UpdatePost(postId, post)
	if err != nil {
		if strings.Contains(err.Error(), "title already exists") {
			resp := utils.ResponseError(nil, err.Error())
			ctx.JSON(http.StatusConflict, resp)
			return
		}
		resp := utils.ResponseError(nil, err.Error())
		ctx.JSON(http.StatusBadGateway, resp)
		return
	}

	resp := utils.ResponseSuccess(updatePost, "Successfully updated post")
	ctx.JSON(http.StatusCreated, resp)
	return
}

func (pc *PostController) ListPost(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "1")
	perPage := ctx.DefaultQuery("per_page", "10")

	intPage, _ := strconv.Atoi(page)
	intPerPage, _ := strconv.Atoi(perPage)

	listPost, err := pc.postService.FindPosts(intPage, intPerPage)

	if err != nil {
		resp := utils.ResponseError(nil, err.Error())
		ctx.JSON(http.StatusBadGateway, resp)
		return
	}

	paginate := utils.PaginationCustom(intPage, intPerPage, len(listPost), listPost)
	resp := utils.ResponseSuccess(paginate, "Successfully get data post")
	ctx.JSON(http.StatusOK, resp)
	return
}

func (pc *PostController) FindPostById(ctx *gin.Context) {
	postId := ctx.Param("postId")

	post, err := pc.postService.FindPostById(postId)

	if err != nil {
		if strings.Contains(err.Error(), "Id exists") {
			resp := utils.ResponseError(nil, err.Error())
			ctx.JSON(http.StatusNotFound, resp)
			return
		}
		resp := utils.ResponseError(nil, err.Error())
		ctx.JSON(http.StatusBadGateway, resp)
		return
	}

	resp := utils.ResponseSuccess(post, "Successfully det detail post")
	ctx.JSON(http.StatusOK, resp)
	return
}

func (pc *PostController) DeletePost(ctx *gin.Context) {
	postId := ctx.Param("postId")
	err := pc.postService.DeletePost(postId)
	if err != nil {
		if strings.Contains(err.Error(), "Id exists") {
			resp := utils.ResponseError(nil, err.Error())
			ctx.JSON(http.StatusNotFound, resp)
			return
		}
		resp := utils.ResponseError(nil, err.Error())
		ctx.JSON(http.StatusBadGateway, resp)
		return
	}

	resp := utils.ResponseSuccess(nil, "Successfully deleted post")
	ctx.JSON(http.StatusOK, resp)
	return
}
