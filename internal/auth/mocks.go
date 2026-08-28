package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TestAuthenticator struct{}

const secret = "test"

var testClaims = jwt.MapClaims{
	"sub": int64(1),
	"exp": time.Now().Add(24 * 7 * 365 * time.Hour).Unix(), // 1 year
	"aud": "test-aud",
	"iss": "test-aud",
}

func (ta *TestAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, testClaims)
	return token.SignedString([]byte(secret))
}
func (ta *TestAuthenticator) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})
}
