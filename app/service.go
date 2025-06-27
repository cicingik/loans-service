package app

import (
	"github.com/cicingik/loans-service/delivery"
)

func initRoutes(
	eng *delivery.HTTPEngine,
	lctrl *delivery.LoansController,
	uctrl *delivery.UsersController,
	actrl *delivery.AssessmentController,
	fctrl *delivery.LoanFundingController,
) error {
	// Default Controller
	eng.Mux.Get("/", delivery.IndexHandler)
	eng.Mux.Get("/v", delivery.VersionHandler)
	eng.Mux.Get("/healthzx", delivery.HealthZX)

	// Resource COontroller
	eng.Mux.Mount("/v1/loan", lctrl.Routes())
	eng.Mux.Mount("/v1/user", uctrl.Routes())
	eng.Mux.Mount("/v1/assessment", actrl.Routes())
	eng.Mux.Mount("/v1/funding", fctrl.Routes())

	return nil
}
