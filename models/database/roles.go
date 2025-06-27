package database

import (
	"context"

	"gorm.io/gorm"
)

type (
	// Roles ...
	Roles struct {
		Model
		Description string `gorm:"Column:description" json:"description"`
	}

	// BundleRoles ...
	BundleRoles struct {
		db *gorm.DB
		t  Roles
	}
)

// TableName ...
func (t *Roles) TableName() string {
	return "users.roles"
}

// InitRoles ...
func InitRoles(ctx context.Context, g *gorm.DB) *BundleRoles {
	return &BundleRoles{
		db: g.WithContext(ctx),
		t:  Roles{},
	}
}
