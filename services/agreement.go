package services

import (
	"fmt"

	"codeberg.org/go-pdf/fpdf"
	"github.com/cicingik/loans-service/models/database"
)

// CreateLenderAgreement ...
func CreateLenderAgreement(data database.LoanFundingDetail) (string, error) {
	fmt.Printf("lender-agreement-%d-%d.pdf \n", data.LoanID, data.LenderID)

	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, fmt.Sprintf("Lender Agreement for Loan %d from User %d", data.LoanID, data.LenderID))

	name := fmt.Sprintf("lender-agreement-%d-%d.pdf", data.LoanID, data.LenderID)

	err := pdf.OutputFileAndClose(name)
	if err != nil {
		return "", err
	}

	return name, nil
}
