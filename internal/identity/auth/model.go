package auth

import (
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/koubae/game-hangar/internal/errs"
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

	UsernameMinLength = 4
	UsernameMaxLength = 20
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

type NewAccountCredential struct {
	Credential string
	AccountID  uuid.UUID
	ProviderID int64
	Secret     string
	SecretType string
	Verified   bool
	VerifiedAt *time.Time
}

func (p *NewAccountCredential) Validate() error {
	if p.Verified && p.VerifiedAt == nil {
		return errs.AccountCredVerifiedAtRequired
	}
	if !p.Verified && p.VerifiedAt != nil {
		return errs.AccountCredVerifiedNilWhenIsFalse
	}
	return nil
}

var usernamePattern = regexp.MustCompile(`^[\pL\pN_-][\pL\pN_-]*$`)

var reservedUsernameNames = map[string]struct{}{
	"admin":     {},
	"moderator": {},
	"mod":       {},
	"support":   {},
	"system":    {},
	"root":      {},
	"null":      {},
	"undefined": {},
	"owner":     {},
	"staff":     {},
	"dev":       {},
	"developer": {},
	"gm":        {},
	"gameadmin": {},
}

func (p *NewAccountCredential) ValidateForTypeUsername() error {
	err := p.Validate()
	if err != nil {
		return err
	}

	p.Credential = strings.TrimSpace(p.Credential)
	if p.Credential == "" {
		return errs.AccountCredCredentialRequired
	}

	length := utf8.RuneCountInString(p.Credential)
	if length < UsernameMinLength {
		return errs.AccountCredCredentialTooShort
	}
	if length > UsernameMaxLength {
		return errs.AccountCredCredentialTooLong
	}

	if !usernamePattern.MatchString(p.Credential) {
		return errs.AccountCredCredentialInvalid
	}
	if _, ok := reservedUsernameNames[strings.ToLower(p.Credential)]; ok {
		return errs.AccountCredCredentialReserved
	}

	return nil
}
