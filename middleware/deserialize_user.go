package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/example/golang-test/config"
	"github.com/example/golang-test/services"
	"github.com/example/golang-test/utils"
	"github.com/gin-gonic/gin"
)

func DeserializeUser(userService services.UserService) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		var token string
		authorizationHeader := ctx.Request.Header.Get("Authorization")
		fields := strings.Fields(authorizationHeader)

		if len(fields) != 0 && fields[0] == "Bearer" {
			token = fields[1]
		}

		if token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You are not logged in"})
			return
		}

		config1, _ := config.LoadConfig1(".")
		sub, err := utils.ValidateToken(token, config1.AccessTokenPrivateKey)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		result, err := userService.FindUserById(fmt.Sprint(sub))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user belonging to this token no logger exists"})
			return
		}
		ctx.Set("currentUser", result)
		ctx.Next()

	}
}
func DeserializeAdmin(adminService services.AdminService) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		var token string
		authorizationHeader := ctx.Request.Header.Get("Authorization")
		fields := strings.Fields(authorizationHeader)

		if len(fields) != 0 && fields[0] == "Bearer" {
			token = fields[1]
		}

		if token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You are not logged in"})
			return
		}

		config1, _ := config.LoadConfig1(".")

		sub, err := utils.ValidateToken(token, config1.AccessTokenPrivateKey)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": err.Error()})
			return
		}

		result, err := adminService.FindAdminByEmail(fmt.Sprint(sub))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "Error Token"})
			return
		}

		ctx.Set("currentUser", result)
		ctx.Next()

	}

}
