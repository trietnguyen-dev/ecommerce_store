package controllers

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/example/golang-test/config"
	"github.com/example/golang-test/models"
	"github.com/example/golang-test/services/products"
	"github.com/example/golang-test/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
	"strconv"
)

type ProductController struct {
	productService products.ProductService
	redis          *redis.Client
	//storageClient  *storage.Client
	config1   *config.Config
	s3Service *s3.S3
}

func NewProductController(productService products.ProductService, redis *redis.Client, config1 *config.Config, s3Service *s3.S3) ProductController {
	return ProductController{productService, redis, config1, s3Service}
}
func (pc *ProductController) CreateProduct(ctx *gin.Context) {
	productInput := ctx.MustGet("product").(models.Product)

	err := pc.productService.CreateProduct(&productInput)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": gin.H{}})
}
func (pc ProductController) DeteleProduct(ctx *gin.Context) {

	productId := ctx.Query("productId")
	err := pc.productService.DeleteProduct(productId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{})
}
func (pc ProductController) GetProductById(ctx *gin.Context) {
	productId := ctx.Query("productId")
	product, err := pc.productService.GetProductById(productId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "product": product})
}

// TODO: edit url params
func (pc ProductController) UploadImage(ctx *gin.Context) {

	file, handler, err := ctx.Request.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

	defer file.Close()

	sizeFile := handler.Size
	buffer := make([]byte, sizeFile) // read file content to buffer

	file.Read(buffer)
	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)
	params := &s3.PutObjectInput{
		Bucket:        aws.String(pc.config1.AwsBucketName),
		Key:           aws.String(handler.Filename),
		Body:          fileBytes,
		ContentLength: aws.Int64(sizeFile),
		ContentType:   aws.String(fileType),
	}
	_, err = pc.s3Service.PutObject(params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}
	// Tạo URL của tệp tin trên S3
	s3URL := "https://s3.amazonaws.com/" + pc.config1.AwsBucketName + "/" + handler.Filename

	ctx.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
		"url":     s3URL,
		"error":   false,
	})
}

func (pc ProductController) UpdateProduct(ctx *gin.Context) {
	productInput := ctx.MustGet("product").(models.Product)
	id := ctx.Query("id")

	err := pc.productService.UpdateProduct(id, &productInput)

	if err != nil {
		ctx.JSON(http.StatusConflict, gin.H{"status": "error", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})

}
func (pc *ProductController) GetListProducts(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		panic(err)
	}

	listProducts, count, err := pc.productService.GetListProducts(int64(page))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	pagination := utils.GetPagination(count)
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "total": count, "pagination": pagination, "Products": listProducts})

}
