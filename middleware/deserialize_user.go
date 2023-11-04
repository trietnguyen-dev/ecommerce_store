package middleware

import (
	"fmt"
	"github.com/example/golang-test/config"
	"github.com/example/golang-test/services/admin"
	"github.com/example/golang-test/services/user"
	"github.com/example/golang-test/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func DeserializeUser(userService user.UserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken, err := ctx.Cookie("access_token")
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Access token not found in cookie"})
			return
		}
		ctx.Request.Header.Set("Authorization", "Bearer "+accessToken)
		config1, _ := config.LoadConfig1(".")

		sub, err := utils.ValidateToken(accessToken, config1.AccessTokenPrivateKey)
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

func DeserializeAdmin(adminService admin.AdminService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Lấy mã thông báo truy cập từ cookie
		accessToken, err := ctx.Cookie("access_token")
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Access token not found in cookie"})
			return
		}
		ctx.Request.Header.Set("Authorization", "Bearer "+accessToken)
		config1, _ := config.LoadConfig1(".")

		sub, err := utils.ValidateToken(accessToken, config1.AccessTokenPrivateKey)
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
