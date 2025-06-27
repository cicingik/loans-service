package delivery

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/cicingik/loans-service/config"
	"github.com/cicingik/loans-service/models/entity"
	"github.com/cicingik/loans-service/pkg/httputils"
	"github.com/cicingik/loans-service/pkg/middleware"
	"github.com/cicingik/loans-service/services"
	"github.com/go-chi/chi"
)

// LoanFundingController ...
type LoanFundingController struct {
	svc *services.LoansService
	cfg *config.AppConfig
}

// NewLoanFundingController ...
func NewLoanFundingController(
	cfg *config.AppConfig,
	svc *services.LoansService,
) *LoanFundingController {
	return &LoanFundingController{
		svc: svc,
		cfg: cfg,
	}
}

// Routes ...
func (fc *LoanFundingController) Routes() http.Handler {
	mux := chi.NewRouter()

	mux.Group(func(r chi.Router) {
		r.Use(middleware.JwtAuthentication(fc.cfg))
		r.Use(middleware.CheckRole("lender"))
		r.Post("/{loan_id}", fc.Create)
	})

	return mux
}

// Create ...
func (fc *LoanFundingController) Create(w http.ResponseWriter, r *http.Request) {
	loanID := chi.URLParam(r, "loan_id")
	uid := r.Context().Value("xUID").(uint64)

	intLoanID, err := strconv.Atoi(loanID)
	if err != nil {
		httputils.ErrorResponseAsJSON(w, httputils.HTTPResponseWrapper{
			HttpCode: http.StatusBadRequest,
			IsError:  true,
			Data:     nil,
			Meta:     "Invalid Loan ID",
		})
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		httputils.ErrorResponseAsJSON(w, httputils.HTTPResponseWrapper{
			HttpCode: http.StatusUnprocessableEntity,
			IsError:  true,
			Data:     nil,
			Meta:     err.Error(),
		})
		return
	}

	var fund entity.FundingLoan

	err = json.Unmarshal(body, &fund)
	if err != nil {
		httputils.ErrorResponseAsJSON(w, httputils.HTTPResponseWrapper{
			HttpCode: http.StatusUnprocessableEntity,
			IsError:  true,
			Data:     nil,
			Meta:     err.Error(),
		})
		return
	}

	fund.LoanID = uint64(intLoanID)
	fund.LenderID = uid

	err = fund.Validate()
	if err != nil {
		httputils.ErrorResponseAsJSON(w, httputils.HTTPResponseWrapper{
			HttpCode: http.StatusBadRequest,
			IsError:  true,
			Data:     nil,
			Meta:     err.Error(),
		})
		return
	}

	funding, err := fc.svc.Funding(fund)
	if err != nil {
		httputils.ErrorResponseAsJSON(w, httputils.HTTPResponseWrapper{
			HttpCode: http.StatusBadRequest,
			IsError:  true,
			Data:     nil,
			Meta:     err.Error(),
		})
		return
	}

	httputils.JsonResponse(w, http.StatusOK, funding, "Funding Sucess")
}
