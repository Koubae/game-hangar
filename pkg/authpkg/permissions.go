package authpkg

import (
	"strings"
	"unicode/utf8"

	"github.com/koubae/game-hangar/pkg/errspkg"
)

const (
	READ     Action = "read"
	WRITE    Action = "write"
	DELETE   Action = "delete"
	WILDCARD Action = "*"
)

type Action string
type Permissions []Scope

type Scope struct {
	Service  string
	Resource string
	Action   Action
}

func ParseScope(scope string) (*Scope, error) {
	parts := strings.Split(scope, ":")

	if len(parts) != 3 {
		return nil, errspkg.AuthPermissionsScopeFormat
	}

	service := strings.ToLower(strings.TrimSpace(parts[0]))
	resource := strings.ToLower(strings.TrimSpace(parts[1]))
	action := Action(strings.ToLower(strings.TrimSpace(parts[2])))

	if utf8.RuneCountInString(service) <= 1 {
		if !IsWildcard(service) {
			return nil, errspkg.AuthPermissionsScopeFormat
		}
	}
	if utf8.RuneCountInString(resource) <= 1 {
		if !IsWildcard(resource) {
			return nil, errspkg.AuthPermissionsScopeFormat
		}
	}
	if utf8.RuneCountInString(string(action)) <= 1 {
		if !action.IsWildcard() {
			return nil, errspkg.AuthPermissionsScopeFormat
		}
	}

	if !action.Valid() {
		return nil, errspkg.AuthPermissionsScopeFormat
	}

	return &Scope{
		Service:  service,
		Resource: resource,
		Action:   action,
	}, nil
}

func (a Action) Valid() bool {
	switch a {
	case READ, WRITE, DELETE, WILDCARD:
		return true
	default:
		return false
	}
}

func (a Action) IsWildcard() bool {
	return a == WILDCARD
}

func IsWildcard(component string) bool {
	return Action(component).IsWildcard()
}
