package account

import (
	"net/mail"
	"strings"
	"time"

	"github.com/koubae/game-hangar/internal/errs"
)

type Account struct {
	ID       string
	Username string
	Email    *string
	Disabled bool
	Created  time.Time
	Updated  time.Time
}

type NewAccount struct {
	Username string
	Email    *string
}

func (p *NewAccount) Validate() error {
	if strings.TrimSpace(p.Username) == "" {
		return errs.UsernameRequired
	}

	if p.Email != nil {
		_, err := mail.ParseAddressList(*p.Email)
		if err != nil {
			return errs.InvalidEmailFormat
		}
	}

	return nil
}
