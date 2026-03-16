package web

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/koubae/game-hangar/pkg/common"
)

type ResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func WriteJSONErrorResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_ = json.NewEncoder(w).Encode(
		&ResponseError{
			Code:    code,
			Message: message,
		},
	)
}

func WriteBusinessErrorResponse(w http.ResponseWriter, err error) {
	code := http.StatusInternalServerError
	message := "unexpected error"
	if businessError, ok := errors.AsType[*common.BusinessError](err); ok {
		code = businessError.HTTPCode
		message = businessError.Message
	}

	WriteJSONErrorResponse(w, code, message)
}
