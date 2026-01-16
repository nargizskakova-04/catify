package dto

type Goal struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	TargetDate  string  `json:"target_date"`
	Progress    float64 `json:"progress"`
	UserID      string  `json:"user_id"`
}
