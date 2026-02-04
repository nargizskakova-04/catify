package domain

type User struct {
	ID       int64  `db:"id"`
	Username string `db:"username"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

type CreateUserInput struct {
	Username string
	Email    string
	Password string
}

type GetUserByIDInput struct {
	ID int64
}

type GetByEmailInput struct {
	Email string
}
