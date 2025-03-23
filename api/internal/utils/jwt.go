package utils

import (
	"fmt"
	"go-api/internal/env"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	secretKey       = []byte(env.GetString("JWT_SECRET_KEY", ""))
	tokenExpiration = time.Duration(env.GetInt("JWT_TOKEN_EXPIRATION", int(time.Hour*10))) // Default to 10 hours
)

type JWTClaims struct {
	UserID int   `json:"user_id"`
	Exp    int64 `json:"exp"`
	Iat    int64 `json:"iat"`
	Nbf    int64 `json:"nbf"`
	jwt.RegisteredClaims
}

// GenerateJWT creates a JWT for a given user ID
func GenerateJWT(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(tokenExpiration).Unix(),
		"iat":     time.Now().Unix(),
		"nbf":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// ValidateJWT checks the validity of a JWT token

func ValidateJWT(tokenString string) (*JWTClaims, error) {
	claims := &JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
