package auth

import (
	"fmt"
	"unicode"
	"unicode/utf8"

	"github.com/koubae/game-hangar/internal/errs"
	"golang.org/x/crypto/bcrypt"
)

type SecretsService struct{}

type SecretsServiceFactory func() *SecretsService

func NewSecretsService() *SecretsService {
	s := &SecretsService{}
	return s
}

func (s *SecretsService) HashSecret(secret string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	if err != nil {
		return "", errs.Wrap(errs.AuthSecretHash, err)
	}

	return string(hash), nil
}

func (s *SecretsService) VerifySecret(hash string, secret string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(secret))
	return err == nil
}

const MinPasswordLengthLowerBound = 3
const MaxPasswordLengthUpperBound = 72

type PasswordValidationRules struct {
	MinLength int
	MaxLength int
	Uppercase bool
	Lowercase bool
	Digits    bool
	Special   bool
}

func NewPasswordValidationRules(
	minLength int,
	maxLength int,
	upperCase bool,
	lowerCase bool,
	digits bool,
	special bool,
) *PasswordValidationRules {
	if minLength < MinPasswordLengthLowerBound {
		minLength = MinPasswordLengthLowerBound
	} else if minLength > MaxPasswordLengthUpperBound {
		minLength = MaxPasswordLengthUpperBound
	}

	if maxLength < MinPasswordLengthLowerBound {
		maxLength = MinPasswordLengthLowerBound
	} else if maxLength > MaxPasswordLengthUpperBound {
		maxLength = MaxPasswordLengthUpperBound
	}

	if minLength > maxLength {
		maxLength = minLength
	}

	return &PasswordValidationRules{
		MinLength: minLength,
		MaxLength: maxLength,
		Uppercase: upperCase,
		Lowercase: lowerCase,
		Digits:    digits,
		Special:   special,
	}
}

type PasswordValidationErrList struct {
	MinLength string
	MaxLength string
	Uppercase string
	Lowercase string
	Digits    string
	Special   string
	hasErr    bool
}

func (e *PasswordValidationErrList) Error() string {
	msg := ""
	if e.MinLength != "" {
		msg += e.MinLength
	}
	if e.MaxLength != "" {
		msg += e.MaxLength
	}
	if e.Uppercase != "" {
		msg += e.Uppercase
	}
	if e.Digits != "" {
		msg += e.Digits
	}
	if e.Special != "" {
		msg += e.Special
	}
	return msg
}

func (s *SecretsService) ValidatePassword(password string, rules PasswordValidationRules) error {
	err := &PasswordValidationErrList{}

	length := utf8.RuneCountInString(password)
	if length < rules.MinLength {
		err.hasErr = true
		err.MinLength = fmt.Sprintf("password below minimum length %d\n", rules.MinLength)
	}
	if length > rules.MaxLength {
		err.hasErr = true
		err.MaxLength = fmt.Sprintf("password above maximum length %d\n", rules.MaxLength)
	}

	var hasUpper bool
	var hasLower bool
	var hasDigit bool
	var hasSpecial bool
	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSpecial = true
		}
	}

	if rules.Uppercase && !hasUpper {
		err.hasErr = true
		err.Uppercase = "at least one uppercase letter is required\n"
	}
	if rules.Lowercase && !hasLower {
		err.hasErr = true
		err.Lowercase = "at least one lowercase letter is required\n"
	}
	if rules.Digits && !hasDigit {
		err.hasErr = true
		err.Digits = "at least one digit is required\n"
	}
	if rules.Special && !hasSpecial {
		err.hasErr = true
		err.Special = "at least one special character is required\n"
	}

	if !err.hasErr {
		return nil
	}
	return err
}
