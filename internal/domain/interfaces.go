package domain

import "context"

type UserRepository interface {
	Create(ctx context.Context, user *User) (int64, error)
	GetUserByID(ctx context.Context, id int64) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
}

type GoalRepository interface {
	Create(ctx context.Context, goal *Goal) (int64, error)
}
