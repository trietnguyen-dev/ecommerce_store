package user

import "github.com/example/golang-test/models"

type UserService interface {
	FindUserById(id string) (*models.DBResponse, error)
	FindUserByEmail(email string) (*models.DBResponse, error)
	UpdateUserById(id string, value *models.UserResponse) error
	ChangePassword(id string, newPassword string) error
}
