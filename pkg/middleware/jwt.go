package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/cicingik/loans-service/config"
	"github.com/cicingik/loans-service/models/entity"
	"github.com/cicingik/loans-service/pkg/httputils"
	"github.com/cicingik/loans-service/repository/auth"
	log "github.com/sirupsen/logrus"
)

var (
	Role   = "xrole"
	UserID = "xUID"
)

func JwtAuthentication(cfg *config.AppConfig) (mw func(http.Handler) http.Handler) {
	mw = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var jwtJson entity.LoginResponse

			jwtAuth, err := auth.NewAuthRepository(cfg)
			if err != nil {
				log.Errorf("Middleware JwtAuth error. Details %v", err)
				httputils.ErrorResponseAsJSON(w, httputils.HTTPResponseWrapper{
					HttpCode: http.StatusUnauthorized,
					IsError:  true,
					Data:     nil,
					Meta:     "Unauthorized",
				})
				return
			}

			token, err := jwtAuth.VerifyToken(r)
			if err != nil {
				log.Errorf("token is not valid: %s", err)
				httputils.ErrorResponseAsJSON(w, httputils.HTTPResponseWrapper{
					HttpCode: http.StatusUnauthorized,
					IsError:  true,
					Data:     nil,
					Meta:     "Unauthorized",
				})

				return
			}

			err = json.Unmarshal([]byte(token), &jwtJson)
			if err != nil {
				log.Errorf("error unmarshalling token: %s", err)
				httputils.ErrorResponseAsJSON(w, httputils.HTTPResponseWrapper{
					HttpCode: http.StatusUnprocessableEntity,
					IsError:  true,
					Data:     nil,
					Meta:     err.Error(),
				})
				return
			}

			ctx := context.WithValue(r.Context(), Role, jwtJson.AccessRoles)
			ctx = context.WithValue(ctx, UserID, jwtJson.UserID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
	return
}
