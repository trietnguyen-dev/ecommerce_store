package routes

import (
	"github.com/example/golang-test/controllers"
	"github.com/example/golang-test/middleware"
	"github.com/example/golang-test/services/admin"
	"github.com/gin-gonic/gin"
)

type AdminRouteController struct {
	adminController controllers.AdminController
}

func NewAdminRouteController(adminController controllers.AdminController) AdminRouteController {
	return AdminRouteController{adminController}
}

func (ac *AdminRouteController) AdminRoutes(rg *gin.RouterGroup, adminService admin.AdminService) {
	router := rg.Group("/admin")

	router.POST("/login", ac.adminController.Login)
	router.Use(middleware.DeserializeAdmin(adminService))
	router.Use(middleware.RequireLoggedIn())

	router.GET("/getlistusers", ac.adminController.GetListUsers)
	router.GET("/getuser", ac.adminController.GetUserById)
	router.PUT("/updateuser", ac.adminController.UpdateUserById)
}
