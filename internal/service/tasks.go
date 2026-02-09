package service

import (
	"catify/internal/domain"
	"context"

	"github.com/rs/zerolog"
)

type TaskRepository interface {
	Create(ctx context.Context, task *domain.Task) (int64, error)
}

type TaskService struct {
	repo   TaskRepository
	logger zerolog.Logger
}

func NewTaskService(repo TaskRepository, logger zerolog.Logger) *TaskService {
	return &TaskService{
		repo:   repo,
		logger: logger,
	}
}
