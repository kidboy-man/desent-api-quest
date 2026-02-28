package services

import (
	"time"

	apperrors "github.com/kidboy-man/8-level-desent/app/errors"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	jwtSecret []byte
}

func NewAuthService(secret string) *AuthService {
	return &AuthService{jwtSecret: []byte(secret)}
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	if username == "" || password == "" {
		return "", apperrors.NewBadRequest("username and password are required")
	}

	claims := jwt.MapClaims{
		"sub": username,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", apperrors.NewInternalServerError("failed to generate token")
	}
	return tokenStr, nil
}

func (s *AuthService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apperrors.NewUnauthorized("unexpected signing method")
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return nil, apperrors.NewUnauthorized("invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, apperrors.NewUnauthorized("invalid token")
	}

	return claims, nil
}
