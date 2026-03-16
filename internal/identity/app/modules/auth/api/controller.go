package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/koubae/game-hangar/pkg/common"
)

type AuthController struct{}

type RegisterByUsernamePayload struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (p *RegisterByUsernamePayload) Validate() error {
	if err := common.DTOSchemaValidation(p); err != nil {
		return errors.New(fmt.Sprintf("invalid payload: %v", err))
	}
	return nil
}

func (controller *AuthController) RegisterByUsername(w http.ResponseWriter, r *http.Request) {
	var payload RegisterByUsernamePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, fmt.Sprintf("invalid json: %v", err), http.StatusBadRequest)
		return
	}
	err := payload.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, "ok RegisterByUsername\n")
}

func (controller *AuthController) LoginByUsername(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, "ok LoginByUsername\n")
}
