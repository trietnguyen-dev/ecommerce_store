package services

import "github.com/example/golang-test/models"

type UserService interface {
	FindUserById(string) (*models.DBResponse, error)
	FindUserByEmail(string) (*models.DBResponse, error)
	UpdateUserById(id string, value *models.UserResponse) error
	//UpdateOne(field string, value interface{}) (*models.DBResponse, error)
}
