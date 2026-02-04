package domain

import "context"

type UserRepository interface {
	Create(ctx context.Context, input CreateUserInput) (int64, error)
	GetUserByID(ctx context.Context, input GetUserByIDInput) (*User, error)
	GetByEmail(ctx context.Context, input GetByEmailInput) (*User, error)
}

type GoalRepository interface {
	Create(ctx context.Context, input CreateGoalInput) (int64, error)
}
