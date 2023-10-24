package services

import (
	"github.com/example/golang-test/config"
	"github.com/example/golang-test/daos"
	"github.com/example/golang-test/models"
)

type AdminServiceImpl struct {
	dao  *daos.DAO
	conf *config.Config
}

func NewAdminService(dao *daos.DAO, conf *config.Config) *AdminServiceImpl {
	return &AdminServiceImpl{dao: dao, conf: conf}
}

func (as *AdminServiceImpl) FindAdminByEmail(email string) (*models.DBResponse, error) {

	user, err := as.dao.FindAdminByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (as *AdminServiceImpl) GetListUsers(page int64) ([]*models.DBResponse, error) {
	users, err := as.dao.GetListUsers(page)
	if err != nil {
		return nil, err
	}
	return users, nil
}
func (as *AdminServiceImpl) FindAdminById(id string) (*models.DBResponse, error) {

	admin, err := as.dao.FindAdminById(id)
	if err != nil {
		return nil, err
	}
	return admin, nil
}
