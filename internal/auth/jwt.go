package auth

import "github.com/golang-jwt/jwt/v5"

type JWTAuthenticator struct {
	secret   string
	audience string
	issuer   string
}

func NewJWTAuthenticator(secret, audience, issuer string) *JWTAuthenticator {
	return &JWTAuthenticator{
		secret:   secret,
		audience: audience,
		issuer:   issuer,
	}
}

func (j *JWTAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(j.secret))
}

func (j *JWTAuthenticator) ValidateToken(tokenString string) (*jwt.Token, error) {
	return nil, nil
}
