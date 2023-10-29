package controllers

import (
	"context"
	"errors"
	"fmt"
	"github.com/example/golang-test/services/admin"
	"github.com/example/golang-test/services/auth"
	"github.com/example/golang-test/services/user"
	"github.com/redis/go-redis/v9"
	"net/http"
	"time"

	"github.com/example/golang-test/config"
	"github.com/example/golang-test/models"
	"github.com/example/golang-test/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthController struct {
	authService  auth.AuthService
	userService  user.UserService
	adminService admin.AdminService
	ctx          context.Context
	goredis      *redis.Client
}

func NewAuthController(authService auth.AuthService, userService user.UserService, adminService admin.AdminService, ctx context.Context, goredis *redis.Client) AuthController {
	return AuthController{authService, userService, adminService, ctx, goredis}
}

// var (
// 	logger, _ = zap.NewProduction()
// )

func (ac *AuthController) SignUpUser(ctx *gin.Context) {
	var userInput *models.SignUpInput

	if err := ctx.ShouldBindJSON(&userInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if userInput.Password != userInput.PasswordConfirm {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Passwords do not match"})
		return
	}

	//err := ac.authService.IsExistUser(userInput)
	//if err != nil {
	//	ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
	//	return
	//}

	err := ac.authService.SignUpUser(userInput)

	if err != nil {
		ctx.JSON(http.StatusConflict, gin.H{"status": "error", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": gin.H{}})
}

func (ac *AuthController) SignInUser(ctx *gin.Context) {
	var credentials *models.SignInInput

	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	res, err := ac.authService.SignInUser(credentials)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or password"})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if err := utils.VerifyPassword(res.Password, credentials.Password); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or Password"})
		return
	}

	config1, _ := config.LoadConfig1(".")

	// Generate Tokens
	accessToken, err := utils.GenerateToken(config1.AccessTokenExpiresIn, res.ID, config1.AccessTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	refreshToken, err := utils.GenerateToken(config1.RefreshTokenExpiresIn, res.ID, config1.RefreshTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	err = ac.goredis.SetXX(ctx, res.ID.Hex(), refreshToken, 10*time.Hour).Err()
	if err != nil {
		panic(err)
	}

	ctx.SetCookie("access_token", accessToken, config1.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refreshToken, config1.RefreshTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config1.AccessTokenMaxAge*60, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": accessToken})
}

func (ac *AuthController) RefreshAccessToken(ctx *gin.Context) {
	message := "could not refresh access token"

	cookie, err := ctx.Cookie("refresh_token")

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": message})
		return
	}

	config1, _ := config.LoadConfig1(".")

	sub, err := utils.ValidateToken(cookie, config1.RefreshTokenPrivateKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	res, err := ac.authService.FindUserById(fmt.Sprint(sub))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user belonging to this token no logger exists"})
		return
	}

	access_token, err := utils.GenerateToken(config1.AccessTokenExpiresIn, res.ID, config1.AccessTokenPrivateKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("access_token", access_token, config1.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config1.AccessTokenMaxAge*60, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": access_token})
}
func (ac *AuthController) LogoutUser(ctx *gin.Context) {

	ctx.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "false", -1, "/", "localhost", false, true)

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
