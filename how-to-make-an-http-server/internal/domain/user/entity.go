package user

import "time"

type User struct {
	ID        int64
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateUserRequest struct {
	Name  string
	Email string
}

type UpdateUserRequest struct {
	Name string
}
