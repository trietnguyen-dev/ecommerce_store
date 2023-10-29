package main

import (
	"context"
	"fmt"
	"github.com/example/golang-test/services/admin"
	"github.com/example/golang-test/services/auth"
	"github.com/example/golang-test/services/user"
	"log"
	"net/http"

	"github.com/example/golang-test/config"
	"github.com/example/golang-test/controllers"
	"github.com/example/golang-test/daos"
	"github.com/example/golang-test/routes"
	"go.uber.org/zap"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var (
	//conf        *config.Config
	server      *gin.Engine
	ctx         context.Context
	redisclient *redis.Client

	userService         user.UserService
	UserController      controllers.UserController
	UserRouteController routes.UserRouteController

	authService         auth.AuthService
	AuthController      controllers.AuthController
	AuthRouteController routes.AuthRouteController

	adminService         admin.AdminService
	AdminController      controllers.AdminController
	AdminRouteController routes.AdminRouteController
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
	redisclient = redis.NewClient(&redis.Options{
		Addr: config1.RedisUri,
	})

	if _, err := redisclient.Ping(ctx).Result(); err != nil {
		panic(err)
	}

	err = redisclient.SetXX(ctx, "test", "zozoozo", 0).Err()
	if err != nil {
		panic(err)
	}

	fmt.Println("Redis client connected successfully...")
	//Auth
	authService = auth.NewAuthService(dao, &config1)
	AuthController = controllers.NewAuthController(authService, userService, adminService, ctx, redisclient)
	AuthRouteController = routes.NewAuthRouteController(AuthController)
	//User
	userService = user.NewUserService(dao, &config1)
	UserController = controllers.NewUserController(userService, redisclient)
	UserRouteController = routes.NewRouteUserController(UserController)

	//Admin
	adminService = admin.NewAdminService(dao, &config1, &userService)
	AdminController = controllers.NewAdminController(adminService, ctx, redisclient, config1)
	AdminRouteController = routes.NewAdminRouteController(AdminController)

	server = gin.Default()
}

func main() {
	config, err := config.LoadConfig1(".")

	if err != nil {
		log.Fatal("Could not load config", err)
	}

	value, err := redisclient.Get(ctx, "test").Result()

	if err == redis.Nil {
		fmt.Println("key: test does not exist")
	} else if err != nil {
		panic(err)
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:8000", "http://localhost:3000"}
	corsConfig.AllowCredentials = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Content-Type", "Authorization"}
	server.Use(cors.New(corsConfig))

	router := server.Group("/api")

	router.GET("/healthchecker", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": value})
	})

	AuthRouteController.AuthRoutes(router, userService)
	UserRouteController.UserRoutes(router, userService)
	AdminRouteController.AdminRoutes(router, adminService)
	log.Fatal(server.Run(":" + config.Port))
}
