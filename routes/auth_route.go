package routes

import (
	"github.com/example/golang-test/controllers"
	"github.com/example/golang-test/middleware"
	"github.com/example/golang-test/services/admin"
	"github.com/example/golang-test/services/user"
	"github.com/gin-gonic/gin"
)

type AuthRouteController struct {
	authController controllers.AuthController
}

func NewAuthRouteController(authController controllers.AuthController) AuthRouteController {
	return AuthRouteController{authController}
}

func (rc *AuthRouteController) AuthRoutes(rg *gin.RouterGroup, userService user.UserService, adminService admin.AdminService) {
	router := rg.Group("/auth")

	router.POST("/register", middleware.ValidateRegisterUser(), rc.authController.SignUpUser)
	router.POST("/login", rc.authController.SignInUser)
	router.GET("/refresh", rc.authController.RefreshAccessToken)
	router.GET("/refreshAdmin", rc.authController.RefreshAccessTokenAdmin)
	router.GET("/logout", middleware.DeserializeUser(userService), rc.authController.LogoutUser)
	router.GET("/logoutAdmin", middleware.DeserializeAdmin(adminService), rc.authController.LogoutAdmin)
}
