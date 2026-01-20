package dto

type CreateGoalRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	TargetDate  string  `json:"target_date"`
	Progress    float64 `json:"progress"`
	UserID      int64   `json:"user_id"`
}
