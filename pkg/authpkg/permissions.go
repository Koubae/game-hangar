package authpkg

import (
	"slices"
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
	Actions  []Action
}

func ParseScope(scope string) (*Scope, error) {
	parts := strings.Split(scope, ":")

	if len(parts) != 3 {
		return nil, errspkg.AuthPermissionsScopeFormat
	}

	service := strings.ToLower(strings.TrimSpace(parts[0]))
	resource := strings.ToLower(strings.TrimSpace(parts[1]))

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

	rawActions := strings.Split(parts[2], ",")

	var actions []Action
	for i, a := range rawActions {
		if i+1 > 4 {
			return nil, errspkg.AuthPermissionsScopeFormat
		}
		normalized := strings.ToLower(strings.TrimSpace(a))
		if normalized == "" {
			continue
		}

		action := Action(normalized)
		if utf8.RuneCountInString(normalized) <= 1 {
			if !action.IsWildcard() {
				return nil, errspkg.AuthPermissionsScopeFormat
			}
			actions = []Action{WILDCARD} // If there is a Wildcard other actions don't make sense
			break
		}
		if !action.Valid() {
			return nil, errspkg.AuthPermissionsScopeFormat
		}

		actions = append(actions, action)
		if len(actions) > 3 {
			return nil, errspkg.AuthPermissionsScopeFormat
		}
	}
	slices.Sort(actions)

	return &Scope{
		Service:  service,
		Resource: resource,
		Actions:  actions,
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
