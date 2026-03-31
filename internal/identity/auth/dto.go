package auth

import "time"

type DTOProvider struct {
	ID          int64     `json:"id"`
	Source      string    `json:"source"`
	Type        string    `json:"type"`
	DisplayName string    `json:"display_name"`
	Category    string    `json:"category"`
	Disabled    bool      `json:"disabled"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
}

type DTOAccountLoggedIn struct {
	AccountID    string `json:"account_id"    binding:"required"`
	Username     string `json:"username"      binding:"required"`
	LoggedCredID int64  `json:"credential_id" binding:"required"`
}

type DTOAccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires"`
}
