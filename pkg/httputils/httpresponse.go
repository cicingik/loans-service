package httputils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	ContentTypeJSON = "application/json; charset=UTF-8"
)

type (
	HTTPResponse interface {
		StatusCode() int
	}
)

func JsonResponse(w http.ResponseWriter, statusCode int, val interface{}, meta interface{}) {
	var success = false
	if statusCode < http.StatusBadRequest || statusCode == http.StatusTeapot {
		success = true
	}

	jsonreponse, err := json.Marshal(HTTPResponseWrapper{
		HttpCode: statusCode,
		IsError:  !success,
		Data:     val,
		Meta:     meta,
	})

	if err != nil {
		ErrorResponseAsJSON(w, HTTPInternalServerError)
	} else {
		w.Header().Set("content-type", ContentTypeJSON)
		w.WriteHeader(statusCode)
		_, _ = w.Write(jsonreponse)
	}
}
func SendPlainError(w http.ResponseWriter, code int, msg string) {
	_, _ = fmt.Fprint(w, msg)
	w.WriteHeader(code)
}
