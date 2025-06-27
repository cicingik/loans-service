package database

import (
	"context"

	"gorm.io/gorm"
)

type (
	// Users ...
	Users struct {
		Model
		Name   string `gorm:"Column:name" json:"name"`
		RoleID uint64 `gorm:"Column:role_id" json:"role_id"`
	}

	// UserWithRole ...
	UserWithRole struct {
		Model
		Role   Roles  `gorm:"foreignKey:role_id" json:"role"`
		Name   string `gorm:"Column:name" json:"name"`
		RoleID uint64 `gorm:"Column:role_id" json:"role_id"`
	}

	// BundleUsers ...
	BundleUsers struct {
		db *gorm.DB
		t  Users
	}
)

// / TableName ...
func (t *Users) TableName() string {
	return "users.users"
}

// TableName ...
func (t *UserWithRole) TableName() string {
	return "users.users"
}

// InitUsers ...
func InitUsers(ctx context.Context, g *gorm.DB) *BundleUsers {
	return &BundleUsers{
		db: g.WithContext(ctx),
		t:  Users{},
	}
}
