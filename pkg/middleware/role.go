package middleware

import (
	"net/http"
	"slices"
	"strings"

	"github.com/cicingik/loans-service/pkg/httputils"
	log "github.com/sirupsen/logrus"
)

func CheckRole(availRole string) (mw func(http.Handler) http.Handler) {
	mw = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			avrl := strings.Split(availRole, ",")
			role := ctx.Value(Role).(string)

			if !slices.Contains(avrl, role) {
				log.Errorf("unauthorized Role with current role %v", role)
				httputils.ErrorResponseAsJSON(w, httputils.HTTPResponseWrapper{
					HttpCode: http.StatusUnauthorized,
					IsError:  true,
					Data:     nil,
					Meta:     "Unauthorized Role",
				})
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
	return
}
