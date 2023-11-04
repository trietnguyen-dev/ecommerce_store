package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/example/golang-test/config"
	"github.com/example/golang-test/controllers"
	"github.com/example/golang-test/daos"
	"github.com/example/golang-test/routes"
	"github.com/example/golang-test/services/admin"
	"github.com/example/golang-test/services/auth"
	"github.com/example/golang-test/services/products"
	"github.com/example/golang-test/services/user"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"log"
)

var (
	//conf        *config.Config
	server      *gin.Engine
	ctx         context.Context
	redisClient *redis.Client
	//storageClient       *storage.Client
	svcAws              *s3.S3
	userService         user.UserService
	userController      controllers.UserController
	userRouteController routes.UserRouteController

	authService         auth.AuthService
	authController      controllers.AuthController
	authRouteController routes.AuthRouteController

	adminService         admin.AdminService
	adminController      controllers.AdminController
	adminRouteController routes.AdminRouteController

	productService         products.ProductService
	productController      controllers.ProductController
	productRouteController routes.ProductRouteController
)

func init() {
	config1, err := config.LoadConfig1(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	ctx = context.TODO()

	dao, err := daos.NewDAO(&config1)
	if err != nil {
		log.Fatal("failed to init daos", zap.Error(err))
	}
	fmt.Println("MongoDB successfully connected...")

	// Connect to Redis
	redisClient = redis.NewClient(&redis.Options{
		Addr:     config1.RedisUri,
		Password: "",
		DB:       0,
	})

	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		panic(err)
	}

	fmt.Println("Redis client connected successfully...")

	////FireBase client
	//opt := option.WithCredentialsFile("serviceAccountKey.json")
	//
	//storageClient, err = storage.NewClient(context.Background(), opt)
	//if err != nil {
	//	log.Fatal(err)
	//}
	creds := credentials.NewStaticCredentials(config1.AwsAccessKey, config1.AwsSecretAccessKey, "")
	_, err = creds.Get()
	if err != nil {
		log.Fatal(err)
	}
	cfg := aws.NewConfig().WithRegion("ap-southeast-2").WithCredentials(creds)
	svcAws = s3.New(session.New(), cfg)
	fmt.Println("AWS client connected successfully...")

	//Auth
	authService = auth.NewAuthService(dao, &config1)
	authController = controllers.NewAuthController(authService, userService, adminService, ctx, *redisClient)
	authRouteController = routes.NewAuthRouteController(authController)
	//User
	userService = user.NewUserService(dao, &config1)
	userController = controllers.NewUserController(userService, redisClient)
	userRouteController = routes.NewRouteUserController(userController)

	//Admin
	adminService = admin.NewAdminService(dao, &config1, &userService)
	adminController = controllers.NewAdminController(adminService, ctx, redisClient, config1)
	adminRouteController = routes.NewAdminRouteController(adminController)

	//product
	productService = products.NewProductService(dao, &config1)
	productController = controllers.NewProductController(productService, redisClient, &config1, svcAws)
	productRouteController = routes.NewRouteProductController(productController)

	server = gin.Default()
}

func main() {
	config, err := config.LoadConfig1(".")

	if err != nil {
		log.Fatal("Could not load config", err)
	}

	//err = redisClient.Set(ctx, "test", "zozoozo", 10*time.Hour).Err()
	//if err != nil {
	//	log.Fatal("Không thể đặt khóa: ", err)
	//}
	//
	//value, err := redisClient.Get(ctx, "test").Result()
	//
	//if err == redis.Nil {
	//	fmt.Println("Khóa 'test' không tồn tại")
	//} else if err != nil {
	//	log.Fatal("Lỗi khi lấy khóa: ", err)
	//}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:8000", "http://localhost:3000"}
	corsConfig.AllowCredentials = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Content-Type", "Authorization"}
	server.Use(cors.New(corsConfig))

	router := server.Group("/api")
	//
	//router.GET("/healthchecker", func(ctx *gin.Context) {
	//	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": value})
	//})

	authRouteController.AuthRoutes(router, userService, adminService)
	userRouteController.UserRoutes(router, userService)
	adminRouteController.AdminRoutes(router, adminService)
	productRouteController.ProductRoutes(router, productService)
	log.Fatal(server.Run(":" + config.Port))
}
