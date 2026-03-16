package model

type Account struct {
	ID       string `json:"id"`
	Username string `json:"username" binding:"required"`
}
