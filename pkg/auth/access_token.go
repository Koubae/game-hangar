package auth

type AccessToken struct {
	Source      string
	Type        string
	AccountID   int64
	Credential  string
	Issuer      string
	Role        string
	AccessToken string
}
