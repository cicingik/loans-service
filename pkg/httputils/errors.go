package httputils

import (
	"encoding/json"
	"net/http"
)

type (
	HTTPResponseWrapper struct {
		HttpCode int         `json:"-"`
		IsError  bool        `json:"error"`
		Data     interface{} `json:"data"`
		Meta     interface{} `json:"meta"`
	}

	PlainMessage struct {
		Message string `json:"response"`
	}
)

var (
	HTTPInternalServerError = HTTPResponseWrapper{
		HttpCode: http.StatusInternalServerError,
		IsError:  true,
		Data:     nil,
		Meta:     http.StatusText(http.StatusInternalServerError),
	}
)

func (h HTTPResponseWrapper) StatusCode() int {
	return h.HttpCode
}

func ErrorResponseAsJSON(w http.ResponseWriter, httpError HTTPResponse) {
	w.Header().Set("Content-Type", ContentTypeJSON)
	w.WriteHeader(httpError.StatusCode())
	_ = json.NewEncoder(w).Encode(httpError)
}
