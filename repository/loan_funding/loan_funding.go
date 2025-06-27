package funding

import (
	"github.com/cicingik/loans-service/config"
	"github.com/cicingik/loans-service/models/database"
	"github.com/cicingik/loans-service/repository/postgre"
	"gorm.io/gorm"
)

type (
	// LoanFundingRepository ...
	LoanFundingRepository struct {
		DB  *postgre.DbEngine
		Cfg *config.AppConfig
	}
)

// NewLoanFundingRepository ...
func NewLoanFundingRepository(cfg *config.AppConfig, db *postgre.DbEngine) (*LoanFundingRepository, error) {
	return &LoanFundingRepository{
		DB:  db,
		Cfg: cfg,
	}, nil
}

// Create ....
func (l *LoanFundingRepository) Create(trx *gorm.DB, data database.LoanFundings) (database.LoanFundings, error) {
	if trx == nil {
		trx = l.DB.G
	}

	err := trx.Create(&data).Error
	if err != nil {
		return database.LoanFundings{}, err
	}

	return data, nil
}

// NoLenderAgreement ...
func (l *LoanFundingRepository) NoLenderAgreement() (*database.LoanFundingDetail, error) {
	var fund database.LoanFundingDetail
	err := l.DB.G.
		Select(`"loans"."loan_fundings".*, l.status as status`).
		Joins("left join loans.loans l on l.id = loan_id").
		Where(`l.status = ?`, "invested").
		Where(`length("loans"."loan_fundings".letter_aggrement) <= 0`).
		First(&fund).Error
	if err != nil {
		return nil, err
	}

	return &fund, err
}

// UpdateLenderAgreemnt ...
func (l *LoanFundingRepository) UpdateLenderAgreemnt(la string, id int) (err error) {
	var loan database.LoanFundings

	udata := map[string]interface{}{
		"letter_aggrement": la,
	}

	return l.DB.G.Model(&loan).
		Where("id", id).
		Updates(udata).Error
}
