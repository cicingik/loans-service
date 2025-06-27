package entity

import (
	"errors"
	"reflect"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// LoanCreate ...
	LoanCreate struct {
		BorrowerID     uint64
		LoanAmmount    float64 `json:"loan_amount"`
		RatePercentage float64 `json:"rate_percentage"`
	}
	// LoanAssessment ...
	LoanAssessment struct {
		CurrentStatus string
		ExecuteAt     time.Time `json:"execute_at"`
		Status        string    `json:"status"`
		EmployeeID    string    `json:"employee_id"`
		Document      string    `json:"document"`
		LoanID        uint64    `json:"loan_id"`
	}

	FundingLoan struct {
		LenderID      uint64
		LoanID        uint64
		FundingAmount int64 `json:"funding_amount"`
	}
)

func (i LoanAssessment) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.Status, validation.Required, validation.In("approved", "disbursed")),
		validation.Field(&i.Status, validation.When(i.CurrentStatus == "proposed", validation.In("approved"))),
		validation.Field(&i.Status, validation.When(i.CurrentStatus == "invested", validation.In("disbursed"))),
		validation.Field(&i.EmployeeID, validation.Required),
		validation.Field(&i.Document, validation.Required),
		validation.Field(&i.ExecuteAt, validation.Required),
	)
}

func (i FundingLoan) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.FundingAmount, validation.Required, validation.Min(1)),
	)
}

func IsValidStatus(input interface{}) error {
	value := reflect.ValueOf(input)

	kind := value.Kind()

	switch kind {
	case reflect.String:

		return nil

	default:
		return errors.New(`interface not string`)
	}
}
