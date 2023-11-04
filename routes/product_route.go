package routes

import (
	"github.com/example/golang-test/controllers"
	"github.com/example/golang-test/middleware"
	"github.com/example/golang-test/services/products"
	"github.com/gin-gonic/gin"
)

type ProductRouteController struct {
	productController controllers.ProductController
}

func NewRouteProductController(productController controllers.ProductController) ProductRouteController {
	return ProductRouteController{productController}
}

func (pc *ProductRouteController) ProductRoutes(rg *gin.RouterGroup, productService products.ProductService) {
	router := rg.Group("product")

	router.POST("/create_product", middleware.ValidateCreateProduct(), pc.productController.CreateProduct)
	router.DELETE("/delete_product", pc.productController.DeteleProduct)
	router.GET("/get_product", pc.productController.GetProductById)
	router.POST("/upload", pc.productController.UploadImage)
	router.PUT("update_product", middleware.ValidateCreateProduct(), pc.productController.UpdateProduct)
	router.GET("/list_products", pc.productController.GetListProducts)
}
