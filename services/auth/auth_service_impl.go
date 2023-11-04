package auth

import (
	"fmt"
	"github.com/example/golang-test/config"
	"github.com/example/golang-test/daos"
	"strings"
	"time"

	"github.com/example/golang-test/models"
	"github.com/example/golang-test/utils"
)

type AuthServiceImpl struct {
	dao  *daos.DAO
	conf *config.Config
}

func NewAuthService(dao *daos.DAO, conf *config.Config) *AuthServiceImpl {
	return &AuthServiceImpl{dao, conf}
}

func (ac *AuthServiceImpl) SignUpUser(user *models.SignUpInput) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.Email = strings.ToLower(user.Email)
	user.Verified = false
	user.Role = "user"
	switch user.Gender {
	case "male":
		user.Gender = "Male"
	case "female":
		user.Gender = "Female"
	default:
		user.Gender = "other"
	}

	hashedPassword, _ := utils.HashPassword(user.Password)
	user.Password = hashedPassword

	err := ac.dao.SignUpUser(user)
	if err != nil {
		return err
	}

	return nil
}

func (ac *AuthServiceImpl) SignInUser(user *models.SignInInput) (*models.DBResponse, error) {
	result, _ := ac.dao.FindUserByEmail(user.Email)
	if result == nil {
		return nil, fmt.Errorf("Invalid user")
	}

	return result, nil
}
func (ac *AuthServiceImpl) FindUserById(id string) (*models.DBResponse, error) {

	user, err := ac.dao.FindUserById(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (ac *AuthServiceImpl) IsExistUser(user *models.SignUpInput) error {
	err := ac.dao.IsExistUser(user)
	if err != nil {
		return err
	}

	return nil
}
func (ac *AuthServiceImpl) FindAdminByEmail(email string) (*models.DBResponse, error) {
	user, err := ac.dao.FindAdminByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}
