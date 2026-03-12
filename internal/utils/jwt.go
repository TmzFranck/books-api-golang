package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	accessSecret  []byte
	refreshSecret []byte
	accessTTL     = 15 * time.Minute
	refreshTTL    = 1 * time.Hour
)

type Claims struct {
	UserId    uint   `json:"user_id"`
	UserEmail string `json:"user_email"`
	Refresh   bool   `json:"refresh,omitempty"`
	jwt.StandardClaims
}

func InitJWT(accessSecretKey, refreshSecretKey string) {
	accessSecret = []byte(accessSecretKey)
	refreshSecret = []byte(refreshSecretKey)
}

func GenerateTokens(userId uint, userEmail string) (string, string, error) {
	// Access Token
	accessClaims := &Claims{
		UserId:    userId,
		UserEmail: userEmail,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(accessTTL).Unix(),
			Issuer:    "BooksApi",
		},
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(accessSecret)
	if err != nil {
		return "", "", fmt.Errorf("could not create access token: %v", err)
	}

	// Refresh Token
	refreshClaims := &Claims{
		UserId:    userId,
		UserEmail: userEmail,
		Refresh:   true,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(refreshTTL).Unix(),
			Issuer:    "BooksApi",
		},
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(refreshSecret)
	if err != nil {
		return "", "", fmt.Errorf("could not create refresh token: %v", err)
	}

	return accessToken, refreshToken, nil
}

func RefreshToken(r string) (string, error) {
	claims, err := ValidateToken(r)
	if !claims.Refresh {
		return "", fmt.Errorf("invalid refresh token")
	}

	// Generate a new access token
	newAccessToken, _, err := GenerateTokens(claims.UserId, claims.UserEmail)
	if err != nil {
		return "", fmt.Errorf("could not generate new access token: %v", err)
	}
	return newAccessToken, nil
}

func ValidateToken(tokenString string) (*Claims, error) {
	// Parse and validate the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return accessSecret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}
