package authpkg

type AccessToken struct {
	Source      string
	Type        string
	AccountID   string
	Credential  string
	Issuer      string
	Role        string
	AccessToken string
}
