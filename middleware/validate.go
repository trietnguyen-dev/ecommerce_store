package middleware

import (
	"github.com/asaskevich/govalidator"
	"github.com/example/golang-test/models"
	"github.com/example/golang-test/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

func ValidateRegisterUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var userInput models.SignUpInput
		if err := ctx.ShouldBindJSON(&userInput); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
			ctx.Abort()
			return
		}

		if govalidator.IsNull(userInput.Email) ||
			govalidator.IsNull(userInput.Name) ||
			govalidator.IsNull(userInput.Password) ||
			govalidator.IsNull(userInput.PhoneNumber) ||
			govalidator.IsNull(userInput.PasswordConfirm) ||
			govalidator.IsNull(string(userInput.Gender)) {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "All fields are required"})
			ctx.Abort()
			return
		}

		if !govalidator.IsEmail(userInput.Email) {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Format is not a valid email"})
			ctx.Abort()
			return
		}
		if utils.IsValidPhoneNumber(userInput.PhoneNumber) == false {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid phone number"})
			ctx.Abort()
			return
		}
		if userInput.Password != userInput.PasswordConfirm {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Passwords do not match"})
			ctx.Abort()
			return
		}
		ctx.Set("userData", userInput)
		ctx.Next()
	}
}
func ValidateUpdateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var userInput models.UserResponse
		if err := ctx.ShouldBindJSON(&userInput); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
			ctx.Abort()
			return
		}

		if govalidator.IsNull(userInput.Email) ||
			govalidator.IsNull(userInput.Name) ||
			govalidator.IsNull(userInput.PhoneNumber) ||
			govalidator.IsNull(string(userInput.Gender)) {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "All fields are required"})
			ctx.Abort()
			return
		}

		if !govalidator.IsEmail(userInput.Email) {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Format is not a valid email"})
			ctx.Abort()
			return
		}
		if utils.IsValidPhoneNumber(userInput.PhoneNumber) == false {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid phone number"})
			ctx.Abort()
			return
		}

		ctx.Set("userData", userInput)
		ctx.Next()
	}
}
func ValidateCreateProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var productInput models.Product
		if err := ctx.ShouldBindJSON(&productInput); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
			return
		}

		v := reflect.ValueOf(productInput)
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			fieldName := t.Field(i).Name
			// Loại trừ kiểm tra cho CreatedAt và Update
			if fieldName == "CreatedAt" || fieldName == "UpdatedAt" {
				continue
			}
			// Kiểm tra xem trường có giá trị zero không
			if reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()) {
				ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": fieldName + " is required"})
				ctx.Abort()
				return
			}
		}

		ctx.Set("product", productInput)
		ctx.Next()
	}
}
