// Package database ...
package database

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type (
	// LoanFundings ...
	LoanFundings struct {
		LetterAggrement string    `gorm:"Column:letter_aggrement" json:"letter_agreement"`
		FundingAt       time.Time `gorm:"Column:funding_at" json:"funding_at"`
		Model
		LenderID      uint64  `gorm:"Column:lender_id" json:"lender_id"`
		LoanID        uint64  `gorm:"Column:loan_id" json:"loan_id"`
		FundingAmount float64 `gorm:"Column:funding_amount" json:"funding_amount"`
	}

	// LoanFundingDetail ...
	LoanFundingDetail struct {
		LetterAggrement string    `gorm:"Column:letter_aggrement" json:"letter_agreement"`
		FundingAt       time.Time `gorm:"Column:funding_at" json:"funding_at"`
		Model
		LenderID      uint64  `gorm:"Column:lender_id" json:"lender_id"`
		LoanID        uint64  `gorm:"Column:loan_id" json:"loan_id"`
		FundingAmount float64 `gorm:"Column:funding_amount" json:"funding_amount"`
		Status        string  `gorm:"Column:status" json:"status"`
	}

	// BundleLoanFundings ...
	BundleLoanFundings struct {
		db *gorm.DB
		t  LoanFundings
	}
)

// TableName ...
func (t *LoanFundings) TableName() string {
	return "loans.loan_fundings"
}

// TableName ...
func (t *LoanFundingDetail) TableName() string {
	return "loans.loan_fundings"
}

// InitLoanFundings ...
func InitLoanFundings(ctx context.Context, g *gorm.DB) *BundleLoanFundings {
	return &BundleLoanFundings{
		db: g.WithContext(ctx),
		t:  LoanFundings{},
	}
}
