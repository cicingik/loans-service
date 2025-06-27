package delivery

import (
	"net/http"

	"github.com/cicingik/loans-service/config"
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

	mux.Group(func(r chi.Router) {
		r.Use(middleware.JwtAuthentication(lc.cfg))
		r.Use(middleware.CheckRole("lender,admin"))
		r.Get("/", lc.Index)
	})

	return mux
}

// Index ...
func (lc *LoansController) Index(w http.ResponseWriter, _ *http.Request) {
	httputils.JsonResponse(w, http.StatusTeapot, nil, struct {
		Message string `json:"message"`
	}{
		Message: http.StatusText(http.StatusOK),
	})
}

// Create ...
func (lc *LoansController) Create(w http.ResponseWriter, _ *http.Request) {
	httputils.JsonResponse(w, http.StatusTeapot, nil, struct {
		Message string `json:"message"`
	}{
		Message: http.StatusText(http.StatusOK),
	})
}
