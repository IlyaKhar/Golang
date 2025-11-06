package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrNoAuthHeader = errors.New("missing Authorization header")
	ErrInvalidToken = errors.New("invalid token")
	ErrUnauthorized = errors.New("unauthorized")
)

type Config struct {
	Secret    []byte
	AccessTTL time.Duration
	Issuer    string
	Audience  string
}

type Claims struct {
	UserID int    `json:"uid"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type Manager struct{ cfg Config }

func NewManager(cfg Config) *Manager {
	if cfg.AccessTTL <= 0 {
		cfg.AccessTTL = 15 * time.Minute
	}
	return &Manager{cfg: cfg}
}

func (m *Manager) GenerateAccessToken(userID int, email, role string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.cfg.Issuer,
			Audience:  jwt.ClaimStrings{m.cfg.Audience},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.cfg.AccessTTL)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(m.cfg.Secret)
}

func (m *Manager) ParseAndValidate(tokenStr string) (*Claims, error) {
	parser := jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	token, err := parser.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return m.cfg.Secret, nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}
	if m.cfg.Issuer != "" && claims.Issuer != m.cfg.Issuer {
		return nil, ErrInvalidToken
	}
	if m.cfg.Audience != "" {
		ok := false
		for _, a := range claims.Audience {
			if a == m.cfg.Audience {
				ok = true
				break
			}
		}
		if !ok {
			return nil, ErrInvalidToken
		}
	}
	return claims, nil
}
