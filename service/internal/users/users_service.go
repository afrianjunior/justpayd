package users

import (
	"context"

	"github.com/afrianjunior/justpayd/internal/pkg"
)

type UserService interface {
	CreateUser(payload *CreateUserRequest) error
	GetUserByEmail(ctx context.Context, email string) (*pkg.User, error)
}

type userService struct {
	userRepository UserRepository
}

func NewUserService(userRepository UserRepository) UserService {
	return &userService{userRepository: userRepository}
}

func (s *userService) CreateUser(payload *CreateUserRequest) error {
	return s.userRepository.CreateUser(payload)
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*pkg.User, error) {
	return s.userRepository.GetUserByEmail(ctx, email)
}
