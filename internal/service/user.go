package service

import (
	"catify/internal/domain"
	"catify/internal/dto"
	"context"

	"github.com/rs/zerolog"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (int64, error)
}

type UserService struct {
	repo   UserRepository
	logger zerolog.Logger
}

func NewUserService(repo UserRepository, logger zerolog.Logger) *UserService {
	return &UserService{
		repo:   repo,
		logger: logger,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req dto.CreateUserRequest) (*domain.User, error) {
	user := &domain.User{
		Username: req.Username,
		Email:    req.Email,
	}
	id, err := s.repo.Create(ctx, user)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to create user in repository")
		return nil, err
	}

	user.ID = id
	return user, nil
}
