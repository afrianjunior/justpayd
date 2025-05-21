package pkg

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// UserContext is the key used to store user information in the request context
type UserContext string

const (
	// UserIDKey is the key used to store the authenticated user ID in the context
	UserIDKey UserContext = "user_id"
	// UserKey is the key used to store the full user object in the context
	UserKey UserContext = "user"
)

// Claims defines the structure for JWT claims with just user_id
// @Schema
type Claims struct {
	// UserID is the unique identifier of the authenticated user
	UserID int `json:"user_id" example:"1"`
	jwt.RegisteredClaims
}

func JWTAuth(config *Config, db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			// Check if the Authorization header has the Bearer prefix
			bearerToken := strings.Split(authHeader, " ")
			if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
				http.Error(w, "Invalid authorization format, expected Bearer {token}", http.StatusUnauthorized)
				return
			}

			fmt.Println("config.JWT.Secret", config.JWT.Secret)
			fmt.Println("bearerToken[1]", bearerToken[1])

			// Parse the JWT token
			token, err := jwt.ParseWithClaims(
				bearerToken[1],
				&Claims{},
				func(token *jwt.Token) (interface{}, error) {
					return []byte(config.JWT.Secret), nil
				},
			)

			fmt.Println("token", token)

			// Handle any errors
			if err != nil {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// Validate the token and extract claims
			if claims, ok := token.Claims.(*Claims); ok && token.Valid {
				// Create context with user ID
				ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)

				// Fetch full user details from database
				var user User
				err := db.QueryRowContext(ctx, "SELECT id, name, email, role, created_at FROM users WHERE id = ?", claims.UserID).
					Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt)

				if err != nil {
					if err == sql.ErrNoRows {
						http.Error(w, "User not found", http.StatusUnauthorized)
						return
					}
					http.Error(w, "Error fetching user details", http.StatusInternalServerError)
					return
				}

				// Add full user to context
				ctx = context.WithValue(ctx, UserKey, &user)

				// Serve with enriched context
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			}
		})
	}
}

// GetUserIDFromContext retrieves the user_id from the request context
//
// @Description Helper function to extract user ID from context
// @Summary Get user ID from context
// @Return int The user ID
// @Return bool Whether the user ID was successfully extracted
func GetUserIDFromContext(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(UserIDKey).(int)
	return userID, ok
}

// GetUserFromContext retrieves the full user object from the request context
//
// @Description Helper function to extract full user details from context
// @Summary Get user from context
// @Return *User The user object
// @Return bool Whether the user was successfully extracted
func GetUserFromContext(ctx context.Context) (*User, bool) {
	user, ok := ctx.Value(UserKey).(*User)
	return user, ok
}

func GenerateJWT(userID int, config *Config) (string, error) {
	expirationTime := time.Now().Add(time.Duration(config.JWT.Expiration) * time.Minute)

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.JWT.Secret))
}

// RequireAuth is a convenience wrapper to ensure a route requires authentication
//
// @Description Convenience wrapper for JWTAuth
// @Summary Require authentication middleware
// @Tags middleware
// @Security BearerAuth
func RequireAuth(config *Config, db *sql.DB) func(http.Handler) http.Handler {
	return JWTAuth(config, db)
}
