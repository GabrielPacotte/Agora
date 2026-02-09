package auth

type AccessToken string

func NewAccessToken(raw string) (AccessToken, error) {
	if raw == "" {
		return "", ErrInvalidToken
	}
	return AccessToken(raw), nil
}

func (t AccessToken) IsZero() bool {
	return t == ""
}
