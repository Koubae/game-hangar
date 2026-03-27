package service

import "errors"

var (
	ErrProviderIsDisabled = errors.New(
		"provider is not enabled",
	)
	ErrGetProvider = errors.New(
		"could not retrieve auth provider",
	)
	ErrCreateCredentialIncorrectProviderType = errors.New(
		"incorrect provider type",
	)
)
