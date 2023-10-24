package services

import "github.com/example/golang-test/models"

type AdminService interface {
	FindAdminByEmail(email string) (*models.DBResponse, error)
	FindAdminById(id string) (*models.DBResponse, error)
	GetListUsers(page int64) ([]*models.DBResponse, error)
}
