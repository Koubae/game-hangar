package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/koubae/game-hangar/pkg/common"
)

// Validator interface for validating a Request Body Payload
type Validator interface {
	Validate() error
}

func LoadJsonBody[T Validator](w http.ResponseWriter, r *http.Request) (*T, bool) {
	var payload T
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		WriteBusinessErrorResponse(
			w, &common.ClientResponseError{
				HTTPCode: http.StatusBadRequest,
				Message:  fmt.Sprintf("invalid json: %v", err.Error()),
			},
		)
		return nil, false
	}
	return &payload, true
}
