package dto

type CreateUserRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
}
