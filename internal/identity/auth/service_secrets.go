package auth

import (
	"fmt"
	"unicode"
	"unicode/utf8"

	"github.com/golang-jwt/jwt/v5"
	"github.com/koubae/game-hangar/pkg/authpkg"
	"github.com/koubae/game-hangar/pkg/common"
	"github.com/koubae/game-hangar/pkg/errspkg"
	"golang.org/x/crypto/bcrypt"
)

const (
	MinPasswordLengthLowerBound = 3
	MaxPasswordLengthUpperBound = 72
	MinPasswordDefault          = 8
	MaxPasswordDefault          = 12
)

var passwordRules *PasswordValidationRules

type SecretsService struct{}

type SecretsServiceFactory func() *SecretsService

func NewSecretsService() *SecretsService {
	s := &SecretsService{}
	return s
}

func (s *SecretsService) HashSecret(secret string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	if err != nil {
		return "", errspkg.Wrap(errspkg.AuthSecretHash, err)
	}

	return string(hash), nil
}

func (s *SecretsService) VerifySecret(hash string, secret string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(secret))
	return err == nil
}

func (s *SecretsService) GenerateJWTAccessToken(
	source string,
	_type string,
	accountID string,
	credential string,
	expire int64,
) (string, error) {
	privateKey := authpkg.GetPrivateKey()
	claims := jwt.MapClaims{
		"sub":   accountID,
		"exp":   expire,
		"iss":   "GameHangar-Identity",
		"role":  "account",
		"scope": "",

		"source":     source,
		"type":       _type,
		"credential": credential,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privateKey)
}

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

func LoadPasswordRulesConfig(envPrefix string) {
	minLength := common.GetEnvInt(envPrefix+"AUTH_PASSWORD_MIN_LENGTH", MinPasswordDefault)
	maxLength := common.GetEnvInt(envPrefix+"AUTH_PASSWORD_MAX_LENGTH", MaxPasswordDefault)
	upperCase := common.GetEnvBool(envPrefix+"AUTH_PASSWORD_UPPERCASE", true)
	lowerCase := common.GetEnvBool(envPrefix+"AUTH_PASSWORD_LOWERCASE", true)
	digits := common.GetEnvBool(envPrefix+"AUTH_PASSWORD_DIGITS", true)
	special := common.GetEnvBool(envPrefix+"AUTH_PASSWORD_SPECIAL", true)

	passwordRules = NewPasswordValidationRules(minLength, maxLength, upperCase, lowerCase, digits, special)
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

func (s *SecretsService) ValidatePasswordDefaultRules(password string) error {
	err := s.ValidatePassword(password, *passwordRules)
	if err != nil {
		return errspkg.Wrap(errspkg.AuthPasswordValidation, err)
	}
	return nil
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
