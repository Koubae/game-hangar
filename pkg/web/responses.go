package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/koubae/game-hangar/pkg/common"
	"go.uber.org/zap"
)

type ResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func WriteJSONResponse[T any](w http.ResponseWriter, code int, body T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(body); err != nil {
		common.GetLogger().Error("Failed to write JSON error response", zap.Error(err))
	}
}

func WriteJSONErrorResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(
		&ResponseError{
			Code:    code,
			Message: message,
		},
	); err != nil {
		common.GetLogger().Error("Failed to write JSON error response", zap.Error(err))
	}
}

func WriteBusinessErrorResponse(w http.ResponseWriter, err error) {
	var code int
	var message string

	if businessError, ok := errors.AsType[*common.ClientResponseError](err); ok {
		code = businessError.HTTPCode
		message = businessError.Message
	} else {

		code = http.StatusInternalServerError
		message = "unexpected error"
		logger := common.GetLogger()
		logger.Error("Unmapped Business error", zap.Error(err))
	}

	WriteJSONErrorResponse(w, code, message)
}

type JSONTime time.Time

func (t JSONTime) MarshalJSON() ([]byte, error) {
	ts := time.Time(t).Format("2006-01-02 15:04:05")
	return json.Marshal(ts)
}
