package server

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

/*
<header>.<payload>.<signature>

eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.
eyJzdWIiOiIxMjMiLCJpc3MiOiJteS1hcGkiLCJhdWQiOiJteS.
hNrgH0pB4B1...  (подпись)

header - {alg: "HS256", typ: "JWT", kid: "key-1"}
payload - {sub: "123", iss: "my-api", aud: "my"}
Стандартные поля:
iss (issuer) - кто выдал токен
aud (audience) - кому предназначен
sub (subject) - кто это
exp - когда истекает время жизни токена ( unix сек)
nbf - не раньше чем
iat - когда был создан
jti - уникальный идентификатор

Signing Method
HMAC using S
*/

type Calims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func NewRegisteredClaims(issuer, audience, subject string, ttl time.Duration, jti string) jwt.RegisteredClaims {
	now := time.Now()
	return jwt.RegisteredClaims{
		Issuer:    issuer,
		Audience:  jwt.ClaimStrings{audience},
		Subject:   subject,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        jti,
	}
}

func NewAccessToken(userID string) (string, error) {
	rc := NewRegisteredClaims(
		"my-api",
		"my",
		userID,
		time.Hour,
		"access-token",
	)

	claims := Calims{
		UserID:           userID,
		RegisteredClaims: rc,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	return token.SignedString([]byte("secret"))
}
