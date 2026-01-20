package service

import (
	"catify/internal/domain"
	"catify/internal/dto"
	"context"

	"github.com/rs/zerolog"
)

type GoalRepository interface {
	Create(ctx context.Context, goal *domain.Goal) (int64, error)
}

type GoalService struct {
	repo   GoalRepository
	logger zerolog.Logger
}

func NewGoalService(repo GoalRepository, logger zerolog.Logger) *GoalService {
	return &GoalService{
		repo:   repo,
		logger: logger,
	}
}

func (s *GoalService) CreateGoal(ctx context.Context, req dto.CreateGoalRequest) (*domain.Goal, error) {
	goal := &domain.Goal{
		UserID:      req.UserID,
		Title:       req.Title,
		Description: req.Description,
	}
	id, err := s.repo.Create(ctx, goal)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to create goal in repository")
		return nil, err
	}

	goal.ID = id
	return goal, nil
}
