package user

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrEmailExitsts = errors.New("email already exists")
	ErrInvalidName  = errors.New("invalid name")
	ErrInvalidEmail = errors.New("invalid email")
)
