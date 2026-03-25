package service

import "errors"

var ErrCreateCredentialIncorrectProviderType = errors.New(
	"incorrect provider type",
)
