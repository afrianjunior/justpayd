package users

import (
	"context"
	"database/sql"

	"github.com/afrianjunior/justpayd/internal/pkg"
)

type UserRepository interface {
	CreateUser(user *CreateUserRequest) error
	GetUserByEmail(ctx context.Context, email string) (*pkg.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(payload *CreateUserRequest) error {
	_, err := r.db.Exec("INSERT INTO users (name, email, role) VALUES (?, ?, ?)", payload.Name, payload.Email, payload.Role)
	return err
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*pkg.User, error) {
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
