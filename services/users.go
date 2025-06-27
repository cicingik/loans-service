package services

import (
	"github.com/cicingik/loans-service/models/entity"
	"github.com/cicingik/loans-service/repository/auth"
	"github.com/cicingik/loans-service/repository/users"
)

// UsersService ...
type UsersService struct {
	arepo *auth.Repository
	urepo *users.Repository
}

// NewUsersService ...
func NewUsersService(
	arepo *auth.Repository,
	urepo *users.Repository,
) (*UsersService, error) {
	return &UsersService{
		arepo: arepo,
		urepo: urepo,
	}, nil
}

// Login ...
func (u *UsersService) Login(userName, pwd string) (*entity.LoginResponse, error) {
	loginData, err := u.urepo.GetLoginUser(userName, pwd)
	if err != nil {
		return nil, err
	}

	token, err := u.arepo.CreateToken(*loginData, "authorize")
	if err != nil {
		return nil, err
	}

	return &entity.LoginResponse{
		AccessToken: token,
		AccessRoles: loginData.UserWithRole.Role.Description,
		UserID:      loginData.UserID,
	}, nil
}
