package v1

import (
	"catify/internal/domain"
	"catify/internal/dto"
	"context"
)

type UserService interface {
	CreateUser(ctx context.Context, req dto.CreateUserRequest) (*domain.User, error)
	GetUserByID(ctx context.Context, id int64) (*domain.User, error)
}

type GoalService interface {
	CreateGoal(ctx context.Context, req dto.CreateGoalRequest) (*domain.Goal, error)
}
