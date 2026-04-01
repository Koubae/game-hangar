package authpkg

import "context"

type AccessToken struct {
	Source      string
	Type        string
	AccountID   string
	Credential  string
	Issuer      string
	Role        string
	Permissions Permissions
	AccessToken string
}

var PermissionEmpty = Permissions{}

func GetAccessToken(ctx context.Context) (*AccessToken, bool) {
	accessToken, ok := ctx.Value(ContextKeyAccessToken).(*AccessToken)
	return accessToken, ok
}

func GetPermissions(ctx context.Context) (Permissions, bool) {
	permissions, ok := ctx.Value(ContextKeyPermissions).(Permissions)
	return permissions, ok
}

func GetPermissionsOrDefault(ctx context.Context) Permissions {
	permissions, ok := GetPermissions(ctx)
	if !ok {
		return PermissionEmpty
	}
	return permissions
}
