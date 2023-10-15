package middleware

import (
	"fmt"
	"grcp-api-client-mongo/config"
	"grcp-api-client-mongo/services"
	"grcp-api-client-mongo/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(userService services.UserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var accessToken string

		cookie, err := ctx.Cookie("access_token")

		authorizationHeader := ctx.Request.Header.Get("Authorization")
		fields := strings.Fields(authorizationHeader)

		if len(fields) != 0 && fields[0] == "Bearer" {
			accessToken = fields[1]
		} else if err != nil {
			accessToken = cookie
		}

		if accessToken == "" {
			resp := utils.ResponseError(nil, "You are not logged in")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, resp)
			return
		}

		config, _ := config.LoadConfig(".")
		sub, err := utils.ValidateToken(accessToken, config.AccessTokenPublicKey)
		if err != nil {
			resp := utils.ResponseError(nil, "Unauthorized")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, resp)
			return
		}

		user, err := userService.FindUserById(fmt.Sprint(sub))
		if err != nil {
			resp := utils.ResponseError(nil, err.Error())
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, resp)
			return
		}

		ctx.Set("currentUser", user)
		ctx.Next()
	}
}
