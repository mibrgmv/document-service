package repository

import (
	"context"

	"github.com/mibrgmv/document-service/internal/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserByLogin(ctx context.Context, login string) (*domain.User, error)
	UserExists(ctx context.Context, login string) (bool, error)
}
