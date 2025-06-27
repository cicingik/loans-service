package database

import (
	"context"

	"gorm.io/gorm"
)

type (
	// LoginData ...
	LoginData struct {
		Model
		UserName string `gorm:"Column:user_name" json:"user_name"`
		Password string `gorm:"Column:password" json:"password"`
		UserID   uint64 `gorm:"Column:user_id" json:"user_id"`
	}

	// LoginDataWithRole ...
	LoginDataWithRole struct {
		Model
		UserName     string       `gorm:"Column:user_name" json:"user_name"`
		Password     string       `gorm:"Column:password" json:"password"`
		UserWithRole UserWithRole `gorm:"foreignkey:user_id" json:"user"`
		UserID       uint64       `gorm:"Column:user_id" json:"user_id"`
	}

	// BundleLoginData ...
	BundleLoginData struct {
		db *gorm.DB
		t  LoginData
	}
)

// TableName ...
func (t *LoginData) TableName() string {
	return "users.login_data"
}

// TableName ...
func (t *LoginDataWithRole) TableName() string {
	return "users.login_data"
}

// InitLoginData ...
func InitLoginData(ctx context.Context, g *gorm.DB) *BundleLoginData {
	return &BundleLoginData{
		db: g.WithContext(ctx),
		t:  LoginData{},
	}
}
