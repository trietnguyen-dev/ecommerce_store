package controllers

import (
	"context"
	"errors"
	"github.com/example/golang-test/config"
	"github.com/example/golang-test/models"
	"github.com/example/golang-test/services/admin"
	"github.com/example/golang-test/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strconv"
	"time"
)

type AdminController struct {
	adminService admin.AdminService
	goredis      *redis.Client
	ctx          context.Context
	config1      config.Config
}

func NewAdminController(adminService admin.AdminService, ctx context.Context, goredis *redis.Client, config1 config.Config) AdminController {
	return AdminController{adminService, goredis, ctx, config1}
}

func (ac *AdminController) Login(ctx *gin.Context) {

	var credentials *models.SignInInput

	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	user, err := ac.adminService.FindAdminByEmail(credentials.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or password"})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if err := utils.VerifyPassword(user.Password, credentials.Password); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or Password"})
		return
	}

	if user.Role != "admin" {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Username or role is not admin"})
		return
	}
	config1, _ := config.LoadConfig1(".")

	// Generate Tokens
	accessToken, err := utils.GenerateToken(config1.AccessTokenExpiresIn, user.Email, config1.AccessTokenPrivateKey)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	refreshToken, err := utils.GenerateToken(config1.RefreshTokenExpiresIn, user.Email, config1.RefreshTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	err = ac.goredis.Set(ctx, user.Email, refreshToken, 7*24*time.Hour).Err()
	if err != nil {
		panic(err)
	}
	ctx.Set("currentUser", user.Email)
	ctx.SetCookie("access_token", accessToken, config1.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refreshToken, config1.RefreshTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config1.AccessTokenMaxAge*60, "/", "localhost", false, false)
	ctx.SetCookie("user_email", user.Email, config1.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": accessToken})
}

func (ac *AdminController) GetListUsers(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		panic(err)
	}

	users, count, err := ac.adminService.GetListUsers(int64(page))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	pagination := utils.GetPagination(count)
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "total": count, "pagination": pagination, "UsersList": users})

}
func (ac AdminController) GetUserById(ctx *gin.Context) {
	id := ctx.Query("id")
	user, err := ac.adminService.GetUserById(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "user": user})

}
func (ac AdminController) UpdateUserById(ctx *gin.Context) {

	userInput := ctx.MustGet("userData").(models.UserResponse)

	//var user *models.UserResponse
	id := ctx.Query("id")

	//if err := ctx.ShouldBindJSON(&user); err != nil {
	//	ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "not a valid user"})
	//	return
	//}

	err := ac.adminService.UpdateUserById(id, &userInput)

	if err != nil {
		ctx.JSON(http.StatusConflict, gin.H{"status": "error", "message": err.Error()})
		return

	}

	ctx.JSON(http.StatusOK, gin.H{})

}
func (ac AdminController) DeleteUserById(ctx *gin.Context) {
	id := ctx.Query("id")
	err := ac.adminService.DeleteUserById(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{})

}
