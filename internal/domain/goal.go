package domain

type Goal struct {
	ID          int64   `db:"id"`
	Title       string  `db:"title"`
	Description string  `db:"description"`
	TargetDate  string  `db:"target_date"`
	Progress    float64 `db:"progress"`
	UserID      int64   `db:"user_id"`
}

type CreateGoalInput struct {
	Title       string
	Description string
	TargetDate  string
	Progress    float64
	UserID      int64
}