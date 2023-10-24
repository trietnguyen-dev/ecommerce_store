package services

import (
	"github.com/example/golang-test/config"
	"github.com/example/golang-test/daos"
	"github.com/example/golang-test/models"
)

type UserServiceImpl struct {
	dao  *daos.DAO
	conf *config.Config
}

func NewUserService(dao *daos.DAO, conf *config.Config) *UserServiceImpl {
	return &UserServiceImpl{dao: dao, conf: conf}
}

func (us *UserServiceImpl) FindUserById(id string) (*models.DBResponse, error) {

	user, err := us.dao.FindUserById(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *UserServiceImpl) FindUserByEmail(email string) (*models.DBResponse, error) {
	user, _ := us.dao.FindUserByEmail(email)

	return user, nil
}
func (us *UserServiceImpl) UpdateUserById(id string, user *models.UserResponse) error {
	err := us.dao.UpdateUserById(id, user)
	if err != nil {
		return err
	}
	return nil
}
