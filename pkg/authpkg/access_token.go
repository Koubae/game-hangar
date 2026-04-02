package authpkg

import "context"

var PermissionEmpty = Permissions{}

const (
	AccountRole      = "account"
	AdminAccountRole = "account_admin"
)

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

func GetAccessToken(ctx context.Context) (*AccessToken, bool) {
	accessToken, ok := ctx.Value(contextKeyAccessToken{}).(*AccessToken)
	return accessToken, ok
}

func GetPermissions(ctx context.Context) (Permissions, bool) {
	permissions, ok := ctx.Value(contextKeyPermissions{}).(Permissions)
	return permissions, ok
}

func GetPermissionsOrDefault(ctx context.Context) Permissions {
	permissions, ok := GetPermissions(ctx)
	if !ok {
		return PermissionEmpty
	}
	return permissions
}

func WithPermissions(ctx context.Context, p Permissions) context.Context {
	return context.WithValue(ctx, contextKeyPermissions{}, p)
}
