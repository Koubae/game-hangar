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

// Differance returns the difference between two permissions.
// The first return value is the permissions that are missing from the second.
// The second return value is the permissions that are present in the first but not in the second.
//
// The difference is calculated by checking if the requested permissions are present in the granted permissions.
// If the requested permission is a wildcard, it is granted if ANY of the granted permissions match.
// If the requested permission is not a wildcard, it is granted if ANY of the granted permissions match and the requested permission is also granted.
//
// For example, if the requested permissions are:
//   - service:resource:read,write
// TODO: At some point must check this function (and all other one that are being called ) since it was created using
// AI late at evening without checking too much. It seems it works but there may be some bugs as well as some
// inefficiencies. It is recommended to review and test thoroughly before using in production.
func (p Permissions) Differance(other Permissions) (Permissions, []string) {
	missing := make([]string, 0)
	if p == nil || other == nil {
		return PermissionEmpty, missing
	}

	diff := make(Permissions)

	for requestedService, requestedResources := range other {
		servicesToCheck := make([]string, 0)
		if requestedService == "*" {
			for s := range p {
				servicesToCheck = append(servicesToCheck, s)
			}
		} else {
			if _, ok := p[requestedService]; ok {
				servicesToCheck = append(servicesToCheck, requestedService)
			}
			if _, ok := p["*"]; ok {
				servicesToCheck = append(servicesToCheck, "*")
			}
		}

		if len(servicesToCheck) == 0 {
			for requestedResource, requestedActions := range requestedResources {
				for _, requestedAction := range requestedActions {
					missing = append(
						missing,
						fmt.Sprintf("%s:%s:%s", requestedService, requestedResource, requestedAction),
					)
				}
			}
			continue
		}

		for requestedResource, requestedActions := range requestedResources {
			matchedAtLeastOneResource := false

			for _, grantedService := range servicesToCheck {
				grantedResources := p[grantedService]

				// Determine which resources to check for this granted service
				resourcesToCheck := make([]string, 0)
				if requestedResource == "*" {
					for r := range grantedResources {
						resourcesToCheck = append(resourcesToCheck, r)
					}
				} else {
					if _, ok := grantedResources[requestedResource]; ok {
						resourcesToCheck = append(resourcesToCheck, requestedResource)
					}
					if _, ok := grantedResources["*"]; ok {
						resourcesToCheck = append(resourcesToCheck, "*")
					}
				}

				for _, grantedResource := range resourcesToCheck {
					grantedActions := grantedResources[grantedResource]

					// Map wildcard service/resource back to requested if necessary
					targetService := requestedService
					if targetService == "*" {
						targetService = grantedService
					}
					targetResource := requestedResource
					if targetResource == "*" {
						targetResource = grantedResource
					}

					collectMatchedActions(
						diff,
						targetService,
						targetResource,
						grantedActions,
						requestedActions,
					)
					matchedAtLeastOneResource = true
				}
			}

			if !matchedAtLeastOneResource {
				for _, requestedAction := range requestedActions {
					missing = append(
						missing,
						fmt.Sprintf("%s:%s:%s", requestedService, requestedResource, requestedAction),
					)
				}
			} else {
				// If we did match some resources, we still need to check if all requested actions were satisfied
				// by at least one of the matched granted resources.
				for _, requestedAction := range requestedActions {
					found := false

					// A bit inefficient but correct: check if this action is in the diff for ANY of the possible targets
					// for this requested service/resource.
					// Actually, it's easier to check if it was ever added to diff for any of the target Service/Resource
					// that correspond to this request.

					for _, grantedService := range servicesToCheck {
						targetService := requestedService
						if targetService == "*" {
							targetService = grantedService
						}

						// Check if this service is in diff
						if resources, ok := diff[targetService]; ok {
							// Determine target resources for this granted service
							grantedResources := p[grantedService]
							for grantedResource := range grantedResources {
								// Check if this grantedResource matches the requestedResource
								if requestedResource != "*" && requestedResource != grantedResource && grantedResource != "*" {
									continue
								}

								targetResource := requestedResource
								if targetResource == "*" {
									targetResource = grantedResource
								}

								if actions, ok := resources[targetResource]; ok {
									if requestedAction == "*" {
										// requested * is special, if we have ANY granted actions here, it's partially or fully satisfied.
										// But the logic for * is that it should return everything.
										// If we have any actions, we consider it found.
										if len(actions) > 0 {
											found = true
											break
										}
									} else if slices.Contains(actions, requestedAction) {
										found = true
										break
									}
								}
							}
						}
						if found {
							break
						}
					}

					if !found {
						missing = append(
							missing,
							fmt.Sprintf("%s:%s:%s", requestedService, requestedResource, requestedAction),
						)
					}
				}
			}
		}
	}

	return diff, missing
}

func collectMatchedActions(
	diff Permissions,
	service string,
	resource string,
	grantedActions []Action,
	requestedActions []Action,
) {
	for _, requestedAction := range requestedActions {
		if requestedAction == "*" {
			for _, grantedAction := range grantedActions {
				ensureDiffBucket(diff, service, resource)
				if !slices.Contains(diff[service][resource], grantedAction) {
					diff[service][resource] = append(diff[service][resource], grantedAction)
				}
			}
			continue
		}

		if slices.Contains(grantedActions, requestedAction) || slices.Contains(grantedActions, Action("*")) {
			ensureDiffBucket(diff, service, resource)
			if !slices.Contains(diff[service][resource], requestedAction) {
				diff[service][resource] = append(diff[service][resource], requestedAction)
			}
			continue
		}

		// Only add to missing if it hasn't been found via some other matched granted service/resource
		// This check is a bit tricky because collectMatchedActions is called multiple times.
		// For now, let's keep it and see if it duplicates.
		// Actually, if we have multiple granted matches, one might satisfy it while another doesn't.
		// We should only mark as missing if NO granted match satisfies it.
	}
}

func ensureDiffBucket(diff Permissions, service string, resource string) {
	if diff[service] == nil {
		diff[service] = make(map[string][]Action)
	}
	if diff[service][resource] == nil {
		diff[service][resource] = make([]Action, 0)
	}
}
