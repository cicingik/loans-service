package delivery

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/cicingik/loans-service/config"
	"github.com/cicingik/loans-service/models/entity"
	"github.com/cicingik/loans-service/pkg/httputils"
	"github.com/cicingik/loans-service/pkg/middleware"
	"github.com/cicingik/loans-service/services"
	"github.com/go-chi/chi"
)

// LoansController ...
type LoansController struct {
	svc *services.LoansService
	cfg *config.AppConfig
}

// NewLoanController ...
func NewLoanController(
	cfg *config.AppConfig,
	svc *services.LoansService,
) *LoansController {
	return &LoansController{
		svc: svc,
		cfg: cfg,
	}
}

// Routes ...
func (lc *LoansController) Routes() http.Handler {
	mux := chi.NewRouter()

	mux.Group(func(r chi.Router) {
		r.Use(middleware.JwtAuthentication(lc.cfg))
		r.Use(middleware.CheckRole("borrower"))
		r.Post("/", lc.Create)
	})

	return mux
}

// Create ...
func (lc *LoansController) Create(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value("xUID").(uint64)
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

	var loan entity.LoanCreate

	err = json.Unmarshal(body, &loan)
	if err != nil {
		httputils.ErrorResponseAsJSON(w, httputils.HTTPResponseWrapper{
			HttpCode: http.StatusUnprocessableEntity,
			IsError:  true,
			Data:     nil,
			Meta:     err.Error(),
		})
		return
	}

	loan.BorrowerID = uid

	err = lc.svc.Create(loan)
	if err != nil {
		httputils.ErrorResponseAsJSON(w, httputils.HTTPResponseWrapper{
			HttpCode: http.StatusInternalServerError,
			IsError:  true,
			Data:     nil,
			Meta:     err.Error(),
		})
		return
	}

	httputils.JsonResponse(w, http.StatusOK, nil, "Create Loan Sucess")
}
