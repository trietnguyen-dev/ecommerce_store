package routes

import (
	"github.com/example/golang-test/controllers"
	"github.com/example/golang-test/middleware"
	"github.com/example/golang-test/services/user"
	"github.com/gin-gonic/gin"
)

type UserRouteController struct {
	userController controllers.UserController
}

func NewRouteUserController(userController controllers.UserController) UserRouteController {
	return UserRouteController{userController}
}

func (uc *UserRouteController) UserRoutes(rg *gin.RouterGroup, userService user.UserService) {
	router := rg.Group("users")
	router.Use(middleware.DeserializeUser(userService))
	router.Use(middleware.RequireLoggedIn())

	router.GET("/getme", uc.userController.GetMe)
	router.PUT("/update", uc.userController.Update)
	router.PUT("/changepassword", uc.userController.ChangePassword)
}
