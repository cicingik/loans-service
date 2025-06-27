// Package users ...
package users

import (
	"github.com/cicingik/loans-service/config"
	"github.com/cicingik/loans-service/models/database"
	"github.com/cicingik/loans-service/repository/postgre"
)

type (
	// Repository ...
	Repository struct {
		DB  *postgre.DbEngine
		Cfg *config.AppConfig
	}
)

// NewUsersRepository ...
func NewUsersRepository(cfg *config.AppConfig, db *postgre.DbEngine) (*Repository, error) {
	return &Repository{
		DB:  db,
		Cfg: cfg,
	}, nil
}

// GetLoginUser ...
func (u *Repository) GetLoginUser(userName, pwd string) (*database.LoginDataWithRole, error) {
	var result database.LoginDataWithRole

	err := u.DB.G.
		Preload("UserWithRole").
		Preload("UserWithRole.Role").
		Where("user_name = ? and password = ?", userName, pwd).
		First(&result).Error
	if err != nil {
		return nil, err
	}

	return &result, nil
}
