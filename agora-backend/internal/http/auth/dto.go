package authhttp

type RegisterRequest struct {
	Login       string `json:"login"`
	DisplayName string `json:"display_name"`
	Password    string `json:"password"`
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AuthJSON struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

type EmptyMeta struct{}

type AuthResponse struct {
	Auth AuthJSON  `json:"auth"`
	Meta EmptyMeta `json:"meta"`
}
