package jwt

import (
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"uid"`
	Role   string `json:"role"`
	// Bisa tambah InstitutionID jika perlu
	jwtv5.RegisteredClaims
}

// GenerateToken membuat token baru (Static helper)
func GenerateToken(claims Claims, secret string, ttl time.Duration) (string, time.Time, error) {
	now := time.Now()
	exp := now.Add(ttl)

	claims.RegisteredClaims = jwtv5.RegisteredClaims{
		Issuer:    "mutabaah-api",
		Subject:   claims.UserID,
		IssuedAt:  jwtv5.NewNumericDate(now),
		ExpiresAt: jwtv5.NewNumericDate(exp),
	}

	token := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(secret))
	return t, exp, err
}

// Tokenizer struct untuk dependency injection (Opsional, tapi kita pertahankan agar main.go tidak error)
type Config struct {
	Secret string
	TTL    time.Duration
	Issuer string
}

type Tokenizer struct{ cfg Config }

func New(cfg Config) *Tokenizer { return &Tokenizer{cfg: cfg} }

// Verify memvalidasi token
func (t *Tokenizer) Verify(tokenStr string) (*Claims, error) {
	parsed, err := jwtv5.ParseWithClaims(tokenStr, &Claims{}, func(token *jwtv5.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtv5.SigningMethodHMAC); !ok {
			return nil, jwtv5.ErrSignatureInvalid
		}
		return []byte(t.cfg.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !parsed.Valid {
		return nil, jwtv5.ErrSignatureInvalid
	}
	claims, ok := parsed.Claims.(*Claims)
	if !ok {
		return nil, jwtv5.ErrInvalidType
	}
	return claims, nil
}
