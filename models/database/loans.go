package database

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type (
	// Loans ...
	Loans struct {
		Model
		ApprovedAt       *time.Time `gorm:"Column:approved_at" json:"approved_at"`
		DisburseAt       *time.Time `gorm:"Column:disburse_at" json:"disburse_at"`
		LetterAggrement  string     `gorm:"Column:letter_aggrement" json:"letter_agreement"`
		VisitDocument    string     `gorm:"Column:visit_document" json:"visit_document"`
		Status           string     `gorm:"Column:status" json:"status"`
		DisburseBy       string     `gorm:"Column:disburse_by" json:"disburse_by"`
		ApproveBy        string     `gorm:"Column:approve_by" json:"approve_by"`
		BorrowerID       uint64     `gorm:"Column:borrower_id" json:"borrower_id"`
		Rate             float64    `gorm:"Column:rate" json:"rate"`
		Principal        float64    `gorm:"Column:principal" json:"principal"`
		RatePercentage   float64    `gorm:"Column:rate_percentage" json:"rate_percentage"`
		FundingRemaining float64    `gorm:"Column:funding_remaining" json:"funding_remaining"`
		ROI              float64    `gorm:"Column:roi" json:"roi"`
	}

	// BundleLoans ...
	BundleLoans struct {
		db *gorm.DB
		t  Loans
	}
)

// TableName ...
func (t *Loans) TableName() string {
	return "loans.loans"
}

// InitLoans ...
func InitLoans(ctx context.Context, g *gorm.DB) *BundleLoans {
	return &BundleLoans{
		db: g.WithContext(ctx),
		t:  Loans{},
	}
}
