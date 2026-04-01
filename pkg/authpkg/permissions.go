package authpkg

import (
	"fmt"
	"slices"
	"sort"
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
type Permissions map[string]map[string][]Action

type Scope struct {
	Service  string
	Resource string
	Actions  []Action
}

func ParsePermissions(scopeRaw string) (Permissions, error) {
	scopes := make([]Scope, 0)
	for _, s := range strings.Split(scopeRaw, "|") {
		scope, err := ParseScope(s)
		if err != nil {
			return nil, err
		} else if scope == nil {
			continue
		}
		scopes = append(scopes, *scope)
	}
	return NewPermissions(scopes), nil
}

func NewPermissions(scopes []Scope) Permissions {
	permissions := make(Permissions)
	for _, scope := range scopes {
		if permissions[scope.Service] == nil {
			permissions[scope.Service] = make(map[string][]Action)
		}

		if permissions.isPermissionInActionList(scope.Service, scope.Resource, WILDCARD) {
			continue
		}

		actions := permissions[scope.Service][scope.Resource]
		if slices.Contains(actions, WILDCARD) {
			permissions[scope.Service][scope.Resource] = []Action{WILDCARD}
			continue
		}

		merged := append(actions, scope.Actions...)
		permissions[scope.Service][scope.Resource] = dedupeActions(merged)

	}

	return permissions
}

func dedupeActions(actions []Action) []Action {
	seen := make(map[Action]bool, len(actions))
	out := make([]Action, 0, len(actions))
	for _, a := range actions {
		if a.IsWildcard() {
			return []Action{WILDCARD}
		}
		if seen[a] {
			continue
		}
		seen[a] = true
		out = append(out, a)
	}

	slices.Sort(out)
	return out
}

func ParseScope(scope string) (*Scope, error) {
	parts := strings.Split(scope, ":")
	partsLen := len(parts)
	if partsLen == 1 {
		return nil, nil
	} else if partsLen != 3 {
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

func (p Permissions) isPermissionInActionList(service string, resource string, action Action) bool {
	serviceP := p[service]
	if serviceP == nil {
		return false
	}
	actions := serviceP[resource]
	if actions == nil {
		return false
	}
	return slices.Contains(actions, action)

}

// IsActionGranted checks if the specified action is permitted for a given service and resource based on the current permissions.
func (p Permissions) IsActionGranted(service string, resource string, action Action) bool {
	if p == nil {
		return false
	}

	allowed := func(actions []Action) bool {
		return slices.Contains(actions, WILDCARD) || slices.Contains(actions, action)
	}

	if resources, ok := p[service]; ok { // Permissions contain resource
		if actions, ok := resources[resource]; ok && allowed(actions) { // Exact Match
			return true
		}
		if actions, ok := resources[string(WILDCARD)]; ok && allowed(actions) { // Resource wildcard
			return true
		}
	}

	if resources, ok := p[string(WILDCARD)]; ok {
		if actions, ok := resources[resource]; ok && allowed(actions) { // service wildcard
			return true
		}
		if actions, ok := resources[string(WILDCARD)]; ok && allowed(actions) {
			return true
		}
	}

	return false

}

func (p Permissions) Scope() string {
	if p == nil {
		return ""
	}

	services := make([]string, 0, len(p))
	for service := range p {
		services = append(services, service)
	}
	sort.Strings(services)

	parts := make([]string, 0)

	for _, service := range services {
		resources := p[service]

		resourceNames := make([]string, 0, len(resources))
		for resource := range resources {
			resourceNames = append(resourceNames, resource)
		}
		sort.Strings(resourceNames)

		for _, resource := range resourceNames {
			actions := resources[resource]
			if len(actions) == 0 {
				continue
			}

			actionNames := make([]string, 0, len(actions))
			for _, action := range actions {
				actionNames = append(actionNames, string(action))
			}
			sort.Strings(actionNames)

			parts = append(parts, fmt.Sprintf("%s:%s:%s", service, resource, strings.Join(actionNames, ",")))
		}
	}

	return strings.Join(parts, "|")
}

func (p Permissions) Differance(other Permissions) (Permissions, []string) {
	missing := make([]string, 0)
	if p == nil || other == nil {
		return PermissionEmpty, missing
	}

	diff := make(Permissions)
	for otherService, otherResources := range other {
		resources, ok := p[otherService]
		if !ok {
			for otherResource, otherActions := range otherResources {
				for _, action := range otherActions {
					missing = append(missing, fmt.Sprintf("%s:%s:%s", otherService, otherResource, action))
				}
			}
			continue
		}

		for otherResource, otherActions := range otherResources {
			actions, ok := resources[otherResource]
			if !ok {
				for _, action := range otherActions {
					missing = append(missing, fmt.Sprintf("%s:%s:%s", otherService, otherResource, action))
				}
				continue
			}

			matchedActions := make([]Action, 0, len(otherActions))
			for _, oAction := range otherActions {
				if slices.Contains(actions, oAction) {
					matchedActions = append(matchedActions, oAction)
					continue
				}
				missing = append(missing, fmt.Sprintf("%s:%s:%s", otherService, otherResource, oAction))
			}

			if len(matchedActions) > 0 {
				if diff[otherService] == nil {
					diff[otherService] = make(map[string][]Action)
				}
				diff[otherService][otherResource] = matchedActions
			}

		}

	}

	return diff, missing
}
