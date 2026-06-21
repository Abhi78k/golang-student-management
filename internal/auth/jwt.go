package auth

import (
	"project-1/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userID string, cfg *config.Config) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(24 * time.Hour),
			),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(
		[]byte(cfg.JWTAccessSecret),
	)
}

func GenerateRefreshToken(userID string, cfg *config.Config) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(360 * time.Hour),
			),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(
		[]byte(cfg.JWTRefreshSecret),
	)
}

func ValidateToken(tokenString string, cfg *config.Config) (*jwt.Token, error) {
	return jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (any, error) {
			return []byte(cfg.JWTAccessSecret), nil
		},
	)
}

func ValidateRefreshToken(
	tokenString string,
	cfg *config.Config,
) (*jwt.Token, error) {

	return jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (any, error) {
			return []byte(cfg.JWTRefreshSecret), nil
		},
	)
}
