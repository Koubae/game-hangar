package auth

import (
	"time"

	"github.com/google/uuid"
)

type ProviderType string

const (
	Username ProviderType = "username"
	Email    ProviderType = "email"
	Device   ProviderType = "device"

	Anonymous ProviderType = "anonymous"
	Guest     ProviderType = "guest"

	Steam      ProviderType = "steam"
	Epic       ProviderType = "epic"
	PSN        ProviderType = "psn"
	Xbox       ProviderType = "xbox"
	Nintendo   ProviderType = "nintendo"
	GPG        ProviderType = "gpg"        // Google Play Games
	GameCenter ProviderType = "gamecenter" // Apple Game Center

	Google   ProviderType = "google"
	Apple    ProviderType = "apple"
	Discord  ProviderType = "discord"
	Facebook ProviderType = "facebook"
)

type Provider struct {
	ID          int64
	Source      string
	Type        string
	DisplayName string
	Category    string
	Disabled    bool
	Created     time.Time
	Updated     time.Time
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
