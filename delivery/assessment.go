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

// AssessmentController ...
type AssessmentController struct {
	svc *services.LoansService
	cfg *config.AppConfig
}

// NewAssessmentController ...
func NewAssessmentController(
	cfg *config.AppConfig,
	svc *services.LoansService,
) *AssessmentController {
	return &AssessmentController{
		svc: svc,
		cfg: cfg,
	}
}

// Routes ...
func (ac *AssessmentController) Routes() http.Handler {
	mux := chi.NewRouter()

	mux.Group(func(r chi.Router) {
		r.Use(middleware.JwtAuthentication(ac.cfg))
		r.Use(middleware.CheckRole("admin"))

		r.Get("/", ac.AssessmentList)
		r.Put("/{loan_id}/{status}", ac.Process)
	})

	return mux
}

// AssessmentList ...
func (ac *AssessmentController) AssessmentList(w http.ResponseWriter, r *http.Request) {
	data, err := ac.svc.UnAssessLoan()
	if err != nil {
		httputils.ErrorResponseAsJSON(w, httputils.HTTPResponseWrapper{
			HttpCode: http.StatusInternalServerError,
			IsError:  true,
			Data:     nil,
			Meta:     err.Error(),
		})
		return
	}

	httputils.JsonResponse(w, http.StatusOK, data, nil)
}

// Process ...
func (ac *AssessmentController) Process(w http.ResponseWriter, r *http.Request) {
	loanID := chi.URLParam(r, "loan_id")
	status := chi.URLParam(r, "status")

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

	loan, err := ac.svc.FindByID(intLoanID)
	if err != nil {
		httputils.ErrorResponseAsJSON(w, httputils.HTTPResponseWrapper{
			HttpCode: http.StatusUnprocessableEntity,
			IsError:  true,
			Data:     nil,
			Meta:     err.Error(),
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

	var req entity.LoanAssessment

	err = json.Unmarshal(body, &req)
	if err != nil {
		httputils.ErrorResponseAsJSON(w, httputils.HTTPResponseWrapper{
			HttpCode: http.StatusUnprocessableEntity,
			IsError:  true,
			Data:     nil,
			Meta:     err.Error(),
		})
		return
	}

	req.Status = status
	req.LoanID = uint64(intLoanID)
	req.CurrentStatus = loan.Status

	err = req.Validate()
	if err != nil {
		httputils.ErrorResponseAsJSON(w, httputils.HTTPResponseWrapper{
			HttpCode: http.StatusBadRequest,
			IsError:  true,
			Data:     nil,
			Meta:     err.Error(),
		})
		return
	}

	_, err = ac.svc.Assess(req)
	if err != nil {
		httputils.ErrorResponseAsJSON(w, httputils.HTTPResponseWrapper{
			HttpCode: http.StatusInternalServerError,
			IsError:  true,
			Data:     nil,
			Meta:     err.Error(),
		})
		return
	}

	httputils.JsonResponse(w, http.StatusOK, nil, "Success")
}
