// Package services ...
package services

import (
	"fmt"
	"time"

	"github.com/cicingik/loans-service/models/database"
	"github.com/cicingik/loans-service/models/entity"
	funding "github.com/cicingik/loans-service/repository/loan_funding"
	"github.com/cicingik/loans-service/repository/loans"
)

// LoansService ...
type LoansService struct {
	lrepo *loans.LoansRepository
	frepo *funding.LoanFundingRepository
}

// NewLoansService ...
func NewLoansService(
	lrepo *loans.LoansRepository,
	frepo *funding.LoanFundingRepository,
) (*LoansService, error) {
	return &LoansService{
		lrepo: lrepo,
		frepo: frepo,
	}, nil
}

// Create ...
func (l *LoansService) Create(data entity.LoanCreate) error {
	return l.lrepo.Create(&database.Loans{
		Status:         "proposed",
		Principal:      data.LoanAmmount,
		BorrowerID:     data.BorrowerID,
		RatePercentage: data.RatePercentage,
		Rate:           data.LoanAmmount * (data.RatePercentage / 100),
		ROI:            data.LoanAmmount * (data.RatePercentage / 100),
	})
}

// UnAssessLoan ...
func (l *LoansService) UnAssessLoan() ([]database.Loans, error) {
	return l.lrepo.UnAssessLoan()
}

// FindByID ...
func (l *LoansService) FindByID(lid int) (database.Loans, error) {
	return l.lrepo.FindByID(lid)
}

// Assess ...
func (l *LoansService) Assess(data entity.LoanAssessment) (*database.Loans, error) {
	if data.Status == "approved" {
		return l.lrepo.Approve(data)
	}

	if data.Status == "disbursed" {
		return l.lrepo.Disburse(data)
	}

	return nil, fmt.Errorf("unvalid status, current status %s", data.Status)
}

// Funding ...
func (l *LoansService) Funding(data entity.FundingLoan) (*database.LoanFundings, error) {
	trx := l.lrepo.DB.G

	err := l.lrepo.Funding(trx, data)
	if err != nil {
		trx.Rollback()
		return nil, err
	}

	funding, err := l.frepo.Create(trx, database.LoanFundings{
		FundingAt:     time.Now(),
		LenderID:      data.LenderID,
		LoanID:        data.LoanID,
		FundingAmount: float64(data.FundingAmount),
	})
	if err != nil {
		trx.Rollback()
		return nil, err
	}

	trx.Commit()

	return &funding, nil
}
