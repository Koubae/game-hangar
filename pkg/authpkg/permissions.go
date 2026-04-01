package authpkg

import (
	"fmt"
	"strings"
)

const (
	READ     Action = "read"
	WRITE    Action = "write"
	DELETE   Action = "delete"
	WILDCARD        = "*"
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
		return nil, fmt.Errorf("invalid scope format: %s", scope)
	}

	return nil, nil
}
