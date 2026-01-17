package v1

import (
	"catify/internal/domain"
	"catify/internal/dto"
	"context"
)

type UserService interface {
	CreateUser(ctx context.Context, req dto.CreateUserRequest) (*domain.User, error)
	GetUser(ctx context.Context, id int64) (*domain.User, error)
}
