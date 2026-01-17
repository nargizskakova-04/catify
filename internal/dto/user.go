package dto

type CreateUserRequest struct {
	Email    string `json:"email"`
	Username string `json:"name"`
	Age      int    `json:"age"`
}
