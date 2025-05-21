package auth

import (
	"context"

	"github.com/afrianjunior/justpayd/internal/pkg"
)

// AuthService provides authentication-related operations
type AuthService interface {
	// VerifyCredentials verifies user credentials and returns a user if valid
	VerifyCredentials(ctx context.Context, email string) (*pkg.User, error)
	// GenerateToken generates a JWT token for the authenticated user
	GenerateToken(userID int) (string, error)
}

// authService is the implementation of AuthService interface
type authService struct {
	authRepository AuthRepository
	config         *pkg.Config
}

// NewAuthService creates a new authentication service
func NewAuthService(authRepository AuthRepository, config *pkg.Config) AuthService {
	return &authService{
		authRepository: authRepository,
		config:         config,
	}
}

// VerifyCredentials verifies user credentials and returns a user if valid
func (s *authService) VerifyCredentials(ctx context.Context, email string) (*pkg.User, error) {
	// In a real application, you would implement password hashing and verification
	// For demo purposes, we're just checking if the user exists
	return s.authRepository.GetUserByEmail(ctx, email)
}

// GenerateToken generates a JWT token for the authenticated user
func (s *authService) GenerateToken(userID int) (string, error) {
	return pkg.GenerateJWT(userID, s.config)
}
