package users_transport_http

import (
	"context"
	"todoapp-try1/internal/core/domain"
)

type UsersHTTPHandler struct {
	userService UserService
}

type UserService interface {
	CreateUser(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)
}
