package controllers

import (
	"github.com/example/golang-test/models"
	"github.com/example/golang-test/services/user"
	"github.com/example/golang-test/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
)

type UserController struct {
	userService user.UserService
	redis       *redis.Client
}

func NewUserController(userService user.UserService, redis *redis.Client) UserController {
	return UserController{userService, redis}
}

func (uc *UserController) GetMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(*models.DBResponse)
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": models.FilteredResponse(currentUser)}})
}

func (uc *UserController) Update(ctx *gin.Context) {
	var user *models.UserResponse
	currentUser := ctx.MustGet("currentUser").(*models.DBResponse)

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	err := uc.userService.UpdateUserById(currentUser.ID.Hex(), user)

	if err != nil {
		ctx.JSON(http.StatusConflict, gin.H{"status": "error", "message": err.Error()})
		return

	}

	ctx.JSON(http.StatusOK, gin.H{})

}

func (uc UserController) ChangePassword(ctx *gin.Context) {
	var userNewPasword *models.PasswordResponse
	currentUser := ctx.MustGet("currentUser").(*models.DBResponse)

	if err := ctx.ShouldBindJSON(&userNewPasword); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if userNewPasword.NewPassword != userNewPasword.ConfirmNewPassword {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Passwords do not match"})
		return
	}

	user, err := uc.userService.FindUserById(currentUser.ID.Hex())

	if err := utils.VerifyPassword(user.Password, userNewPasword.CurrentPassword); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or Password"})
		return
	}

	err = uc.userService.ChangePassword(currentUser.ID.Hex(), userNewPasword.NewPassword)

	if err != nil {
		ctx.JSON(http.StatusConflict, gin.H{"status": "error", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
