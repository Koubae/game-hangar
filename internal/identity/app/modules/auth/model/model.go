package model

import (
	"time"

	"github.com/google/uuid"
)

type Provider struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	Category    string    `json:"category"`
	Disabled    bool      `json:"disabled"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
}

type AccountCredential struct {
	ID         int64
	Credential string
	AccountID  uuid.UUID
	ProviderID int64

	Secret     string
	SecretType string

	Verified   bool
	VerifiedAt *time.Time

	Disabled   bool
	DisabledAt *time.Time

	Created time.Time
	Updated time.Time
}
