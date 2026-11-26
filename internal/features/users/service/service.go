package usersservice

import (
	"context"
	"todoapp-try1/internal/core/domain"
)

type UsersService struct {
	repository Repository
}

type Repository interface {
	CreateUser(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)
}
