package loans

import (
	"fmt"

	"github.com/cicingik/loans-service/config"
	"github.com/cicingik/loans-service/models/database"
	"github.com/cicingik/loans-service/models/entity"
	"github.com/cicingik/loans-service/repository/postgre"
	"gorm.io/gorm"
)

type (
	LoansRepository struct {
		DB  *postgre.DbEngine
		Cfg *config.AppConfig
	}
)

func NewLoanRepository(cfg *config.AppConfig, db *postgre.DbEngine) (*LoansRepository, error) {
	return &LoansRepository{
		DB:  db,
		Cfg: cfg,
	}, nil
}

func (l *LoansRepository) Create(data *database.Loans) error {
	err := l.DB.G.Create(data).Error
	if err != nil {
		return err
	}

	return nil
}

func (l *LoansRepository) FindByID(uid int) (loan database.Loans, err error) {
	err = l.DB.G.Model(database.Loans{}).First(&loan, uid).Error
	if err != nil {
		return loan, err
	}

	return loan, err
}

func (l *LoansRepository) UnAssessLoan() ([]database.Loans, error) {
	var data []database.Loans

	tx := l.DB.G.
		Where(`status = ?`, "proposed").
		Find(&data)

	err := tx.Error
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (l *LoansRepository) Approve(data entity.LoanAssessment) (*database.Loans, error) {
	udata := map[string]interface{}{
		"status":         data.Status,
		"visit_document": data.Document,
		"approve_by":     data.EmployeeID,
		"approve_at":     data.ExecuteAt,
	}

	var loan database.Loans
	err := l.DB.G.First(&loan, data.LoanID).
		Where("status = ?", "proposed").
		Updates(udata).
		Error
	if err != nil {
		return nil, err
	}

	return &loan, err
}

func (l *LoansRepository) Disburse(data entity.LoanAssessment) (*database.Loans, error) {
	udata := map[string]interface{}{
		"status":           data.Status,
		"letter_aggrement": data.Document,
		"disburse_by":      data.EmployeeID,
		"disburse_at":      data.ExecuteAt,
	}

	var loan database.Loans
	err := l.DB.G.First(&loan, data.LoanID).
		Where("status = ?", "invested").
		Updates(udata).
		Error
	if err != nil {
		return nil, err
	}

	return &loan, err
}

func (l *LoansRepository) Funding(trx *gorm.DB, data entity.FundingLoan) (err error) {
	var loan database.Loans

	if trx == nil {
		trx = l.DB.G
	}

	trx = trx.Model(&loan).
		Where("id", data.LoanID).
		Where("status", "approved").
		Where("funding_remaining >= ?", data.FundingAmount).
		UpdateColumn(
			"funding_remaining", gorm.Expr("funding_remaining - ?", data.FundingAmount),
		)

	if trx.RowsAffected < 1 {
		return fmt.Errorf("failed to funding loan id %d", data.LoanID)
	}

	return trx.Error
}

// Invested ...
func (l *LoansRepository) Invested() (*database.Loans, error) {
	udata := map[string]interface{}{
		"status": "invested",
	}

	var loan database.Loans
	err := l.DB.G.Find(&loan).
		Where("status = ?", "approved").
		Where("funding_remaining = ?", 0).
		Updates(udata).
		Error
	if err != nil {
		return nil, err
	}

	return &loan, err
}
