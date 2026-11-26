package usersservice

import (
	"context"
	"fmt"
	"todoapp-try1/internal/core/domain"
)

func (s *UsersService) CreateUser(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	// 1. validate userDomain
	if err := user.Validate(); err != nil {
		return domain.User{}, fmt.Errorf("validate user: %w", err)
	}

	// 2. if validate ok -> s.repository.CreateUser(..)
	// return created user
}
