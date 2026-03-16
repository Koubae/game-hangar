package api

import (
	"io"
	"net/http"
)

type AuthController struct{}

func (controller *AuthController) RegisterByUsername(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, "ok RegisterByUsername\n")
}

func (controller *AuthController) LoginByUsername(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, "ok LoginByUsername\n")
}
