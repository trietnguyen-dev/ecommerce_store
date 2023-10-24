package routes

import (
	"github.com/example/golang-test/controllers"
	"github.com/example/golang-test/middleware"
	"github.com/example/golang-test/services"
	"github.com/gin-gonic/gin"
)

type AdminRouteController struct {
	adminController controllers.AdminController
}

func NewAdminRouteController(adminController controllers.AdminController) AdminRouteController {
	return AdminRouteController{adminController}
}

func (ac *AdminRouteController) AdminRoutes(rg *gin.RouterGroup, adminService services.AdminService) {
	router := rg.Group("/admin")

	router.GET("/login", ac.adminController.Login)
	router.Use(middleware.DeserializeAdmin(adminService))
	router.Use(middleware.RequireLoggedIn())

	router.GET("/getlistusers", ac.adminController.GetListUsers)
}
