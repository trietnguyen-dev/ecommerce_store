package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireLoggedIn() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookie, _ := ctx.Cookie("logged_in")
		if cookie == "" {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "Invalid login"})
			return
		}
		ctx.Next()
	}
}
