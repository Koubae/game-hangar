package model

import (
	"time"

	"github.com/google/uuid"
)

type ProviderType string

const (
	// Managed
	Username ProviderType = "username"
	Email    ProviderType = "email"
	Device   ProviderType = "device"
	// Anonymous
	Anonymous ProviderType = "anonymous"
	Guest     ProviderType = "guest"
	// Platform
	Steam      ProviderType = "steam"
	Epic       ProviderType = "epic"
	PSN        ProviderType = "psn"
	Xbox       ProviderType = "xbox"
	Nintendo   ProviderType = "nintendo"
	GPG        ProviderType = "gpg"        // Google Play Games
	GameCenter ProviderType = "gamecenter" // Apple Game Center
	// Social
	Google   ProviderType = "google"
	Apple    ProviderType = "apple"
	Discord  ProviderType = "discord"
	Facebook ProviderType = "facebook"
)

type Provider struct {
	ID          int       `json:"id"`
	Source      string    `json:"source"`
	Type        string    `json:"type"`
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
