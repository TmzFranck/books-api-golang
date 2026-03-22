package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	accessSecret  []byte             // secret key for access tokens
	refreshSecret []byte             // secret key for refresh tokens
	accessTTL     = 15 * time.Minute // Time to live for access tokens
	refreshTTL    = 1 * time.Hour    // Time to live for refresh tokens
)

// Claims represents the JWT claims for access and refresh tokens
type Claims struct {
	UserId    uint   `json:"user_id"`
	UserEmail string `json:"user_email"`
	Refresh   bool   `json:"refresh,omitempty"`
	jwt.RegisteredClaims
}

// UrlSafeToken represents a URL-safe JWT token for email verification
type UrlSafeToken struct {
	Usermail string `json:"usermail"`
	jwt.RegisteredClaims
}

// InitJWT initializes the JWT secret keys
func InitJWT(accessSecretKey, refreshSecretKey string) {
	accessSecret = []byte(accessSecretKey)
	refreshSecret = []byte(refreshSecretKey)
}

// GenerateTokens generates a new access and refresh token pair for the given user
func GenerateTokens(userId uint, userEmail string) (string, string, error) {
	// Access Token
	accessClaims := &Claims{
		UserId:    userId,
		UserEmail: userEmail,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTTL)),
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
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTTL)),
			Issuer:    "BooksApi",
		},
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(refreshSecret)
	if err != nil {
		return "", "", fmt.Errorf("could not create refresh token: %v", err)
	}

	return accessToken, refreshToken, nil
}

// RefreshToken refreshes the access token using a valid refresh token and returns the new access token
func RefreshToken(r string) (string, error) {
	claims, err := ValidateToken(r)
	if err != nil || !claims.Refresh {
		return "", fmt.Errorf("invalid refresh token")
	}

	// Generate a new access token
	newAccessToken, _, err := GenerateTokens(claims.UserId, claims.UserEmail)
	if err != nil {
		return "", fmt.Errorf("could not generate new access token: %v", err)
	}
	return newAccessToken, nil
}

// ValidateToken validates the given token string and returns the claims if valid
func ValidateToken(tokenString string) (*Claims, error) {
	// Parse and validate the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
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

// GenerateURLSafeToken generates a URL-safe JWT token for email verification
func GenerateURLSafeToken(usermail string) (string, error) {
	token := &UrlSafeToken{
		Usermail: usermail,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			Issuer:    "BooksApi",
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, token).SignedString(accessSecret)
}

// ValidateURLSafeToken validates the given URL-safe JWT token and returns the claims if valid
func ValidateURLSafeToken(tokenString string) (*UrlSafeToken, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UrlSafeToken{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return accessSecret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*UrlSafeToken)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}
