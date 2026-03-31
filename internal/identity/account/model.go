package account

import "time"

type Account struct {
	ID       string
	Username string
	Email    *string
	Disabled bool
	Created  time.Time
	Updated  time.Time
}
