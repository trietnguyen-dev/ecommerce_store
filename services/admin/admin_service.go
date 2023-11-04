package admin

import "github.com/example/golang-test/models"

type AdminService interface {
	FindAdminByEmail(email string) (*models.DBResponse, error)
	FindAdminById(id string) (*models.DBResponse, error)
	GetListUsers(page int64) ([]*models.DBResponse, int64, error)
	GetUserById(id string) (*models.UserResponse, error)
	UpdateUserById(id string, user *models.UserResponse) error
	DeleteUserById(id string) error
}
