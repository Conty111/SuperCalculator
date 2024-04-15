package clierrs

import "errors"

var (
	ErrTaskNotFound       = errors.New("task not found")
	ErrTaskAlreadyCreated = errors.New("task already created")

	ErrTokenExpired                  = errors.New("token time to live is expired")
	ErrAuthTokenWasNotProvided       = errors.New("auth token required")
	ErrInvalidCredentials            = errors.New("incorrect login or password")
	ErrInvalidAuthToken              = errors.New("invalid auth token")
)
