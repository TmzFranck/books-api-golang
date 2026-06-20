package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/TmzFranck/books-api-golang/internal/utils"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type contextKey string

type HandlerFunc func(http.Handler) http.Handler

const (
	UserIDKey    contextKey = "user_id"
	UserEmailKey contextKey = "user_email"
	UserTokenKey contextKey = "user_token"
)

// NewAuthMiddleware creates a new authentication middleware that validates JWT tokens from the Authorization header
func NewAuthMiddleware(redisClient *redis.Client, log *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := log.WithField("module", "AuthMiddleware")
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				logger.Warn("missing authorization header")
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				logger.Warn("invalid authorization header format")
				http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]

			claims, err := utils.ValidateToken(tokenString)
			if err != nil {
				logger.Warn("invalid or expired token")
				http.Error(w, "invalid or expired token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserId)
			ctx = context.WithValue(ctx, UserEmailKey, claims.UserEmail)
			ctx = context.WithValue(ctx, UserTokenKey, tokenString)

			isBlacklisted, err := utils.GetJwtBlacklist(ctx, redisClient)
			if isBlacklisted {
				logger.Errorf("blacklisted token: %v", err)
				http.Error(w, "invalid or expired token", http.StatusUnauthorized)
				return
			}

			if err != nil {
				logger.Errorf("Invalid or expired token %s", tokenString)
				http.Error(w, "invalid or expired token", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID returns the user ID from the context
func GetUserID(ctx context.Context) uint {
	userID, _ := ctx.Value(UserIDKey).(uint)
	return userID
}

// GetUserEmail returns the user email from the context
func GetUserEmail(ctx context.Context) string {
	userEmail, _ := ctx.Value(UserEmailKey).(string)
	return userEmail
}

// GetUserToken returns the user token from the context
func GetUserToken(ctx context.Context) string {
	token, _ := ctx.Value(UserTokenKey).(string)
	return token
}
