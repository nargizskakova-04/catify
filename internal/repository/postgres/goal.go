package postgres

import (
	"catify/internal/domain"
	"context"
	"database/sql"
)

type GoalRepository struct {
	db *sql.DB
}

func NewGoalRepository(db *sql.DB) *GoalRepository {
	return &GoalRepository{db: db}
}

func (r *GoalRepository) Create(ctx context.Context, goal *domain.Goal) (int64, error) {
	var goalID int64
	query := `INSERT INTO goals (title, description, target_date, progress, user_id) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, goal.Title, goal.Description, goal.TargetDate, goal.Progress, goal.UserID).Scan(&goalID)
	if err != nil {
		return 0, err
	}
	return goalID, nil
}
