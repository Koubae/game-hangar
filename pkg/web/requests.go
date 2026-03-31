package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/koubae/game-hangar/pkg/common"
)

func LoadJsonBodyOrBadRequestResponse[T any](w http.ResponseWriter, r *http.Request) (*T, bool) {
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
