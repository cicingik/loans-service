package delivery

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/cicingik/loans-service/config"
	"github.com/cicingik/loans-service/models/database"
	"github.com/cicingik/loans-service/pkg/httputils"
	"github.com/cicingik/loans-service/services"
	"github.com/go-chi/chi"
)

// UsersController ...
type UsersController struct {
	usvc *services.UsersService
	cfg  *config.AppConfig
}

// NewUsersController ...
func NewUsersController(
	cfg *config.AppConfig,
	usvc *services.UsersService,
) *UsersController {
	return &UsersController{
		usvc: usvc,
		cfg:  cfg,
	}
}

// Routes ...
func (uc *UsersController) Routes() http.Handler {
	mux := chi.NewRouter()

	mux.Group(func(r chi.Router) {
		r.Post("/login", uc.Login)
	})

	return mux
}

// Login ...
func (uc *UsersController) Login(w http.ResponseWriter, r *http.Request) {
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

	var loginData database.LoginData

	err = json.Unmarshal(body, &loginData)
	if err != nil {
		httputils.ErrorResponseAsJSON(w, httputils.HTTPResponseWrapper{
			HttpCode: http.StatusUnprocessableEntity,
			IsError:  true,
			Data:     nil,
			Meta:     err.Error(),
		})
		return
	}

	data, err := uc.usvc.Login(loginData.UserName, loginData.Password)
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
