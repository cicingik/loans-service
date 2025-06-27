// Package delivery ...
package delivery

import (
	"net/http"

	"github.com/cicingik/loans-service/config"
	"github.com/cicingik/loans-service/pkg/httputils"
)

// IndexHandler ...
func IndexHandler(w http.ResponseWriter, _ *http.Request) {
	httputils.JsonResponse(w, http.StatusTeapot, nil, struct {
		Message string `json:"message"`
	}{
		Message: http.StatusText(http.StatusTeapot),
	})
}

// VersionHandler ...
func VersionHandler(w http.ResponseWriter, _ *http.Request) {
	httputils.JsonResponse(w, http.StatusOK, nil, struct {
		Version    string `json:"version"`
		CommitHash string `json:"commit_hash"`
		CommitMsg  string `json:"commit_msg"`
	}{
		Version:    config.AppVersion,
		CommitHash: config.CommitHash,
		CommitMsg:  config.CommitMsg,
	})
}

// HealthZX ...
func HealthZX(w http.ResponseWriter, _ *http.Request) {
	httputils.JsonResponse(w, http.StatusOK, nil, struct {
		Message string `json:"message"`
	}{
		Message: `I am Healthy`,
	})
}
