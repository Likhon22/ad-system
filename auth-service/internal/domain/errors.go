package domain

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserSuspended      = errors.New("user is suspended")
	ErrInvalidRole        = errors.New("invalid role")
	ErrOAuthUser          = errors.New("account uses OAuth login")
)
