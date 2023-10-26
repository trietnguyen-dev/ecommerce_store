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

func (ac *AuthServiceImpl) SignUpUser(user *models.SignUpInput) (*models.DBResponse, error) {
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	user.Email = strings.ToLower(user.Email)
	user.PasswordConfirm = ""
	user.Verified = false
	user.Role = "user"

	hashedPassword, _ := utils.HashPassword(user.Password)
	user.Password = hashedPassword

	res, err := ac.dao.SignUpUser(user)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ac *AuthServiceImpl) SignInUser(user *models.SignInInput) (*models.DBResponse, error) {
	result, _ := ac.dao.FindUserByEmail(user.Email)
	if result == nil {
		return nil, fmt.Errorf("Invalid user")
	}

	return result, nil
}
