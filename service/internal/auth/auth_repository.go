package auth

import (
	"context"
	"database/sql"

	"github.com/afrianjunior/justpayd/internal/pkg"
)

// AuthRepository provides authentication-related database operations
type AuthRepository interface {
	// GetUserByEmail retrieves a user by email for authentication
	GetUserByEmail(ctx context.Context, email string) (*pkg.User, error)
}

// authRepository is the implementation of AuthRepository interface
type authRepository struct {
	db *sql.DB
}

// NewAuthRepository creates a new authentication repository
func NewAuthRepository(db *sql.DB) AuthRepository {
	return &authRepository{db: db}
}

// GetUserByEmail retrieves a user by email for authentication
func (r *authRepository) GetUserByEmail(ctx context.Context, email string) (*pkg.User, error) {
	var user pkg.User
	err := r.db.QueryRowContext(ctx, "SELECT id, name, email, role, created_at FROM users WHERE email = ?", email).
		Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}
