package auth

import "github.com/example/golang-test/models"

type AuthService interface {
	SignUpUser(*models.SignUpInput) error
	IsExistUser(*models.SignUpInput) error
	SignInUser(*models.SignInInput) (*models.DBResponse, error)
	FindUserById(id string) (*models.DBResponse, error)
}
