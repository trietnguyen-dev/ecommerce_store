package admin

import (
	"github.com/example/golang-test/config"
	"github.com/example/golang-test/daos"
	"github.com/example/golang-test/models"
	"github.com/example/golang-test/services/user"
	"time"
)

type AdminServiceImpl struct {
	dao         *daos.DAO
	conf        *config.Config
	userService *user.UserService
}

func NewAdminService(dao *daos.DAO, conf *config.Config, userService *user.UserService) *AdminServiceImpl {
	return &AdminServiceImpl{dao: dao, conf: conf, userService: userService}
}

func (as *AdminServiceImpl) FindAdminByEmail(email string) (*models.DBResponse, error) {

	user, err := as.dao.FindAdminByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (as *AdminServiceImpl) GetListUsers(page int64) ([]*models.DBResponse, int64, error) {
	users, count, err := as.dao.GetListUsers(page)
	if err != nil {
		return nil, 0, err
	}
	return users, count, nil
}
func (as *AdminServiceImpl) FindAdminById(id string) (*models.DBResponse, error) {

	admin, err := as.dao.FindAdminById(id)
	if err != nil {
		return nil, err
	}
	return admin, nil
}
func (as AdminServiceImpl) GetUserById(id string) (*models.UserResponse, error) {
	user, err := as.dao.GetUserById(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (as AdminServiceImpl) UpdateUserById(id string, user *models.UserResponse) error {
	user.UpdatedAt = time.Now()
	err := as.dao.UpdateUserById(id, user)
	if err != nil {
		return err
	}
	return nil
}
func (as AdminServiceImpl) DeleteUserById(id string) error {
	err := as.dao.DeleteUserById(id)
	if err != nil {
		return err
	}
	return nil
}
