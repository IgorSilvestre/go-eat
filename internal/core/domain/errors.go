package domain

import "errors"

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrEmailAlreadyExists   = errors.New("email already exists")
	ErrClerkIDAlreadyExists = errors.New("clerk_id already exists")
	ErrUserAlreadyExists    = errors.New("user with this email or clerk_id already exists")
)
