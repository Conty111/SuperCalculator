package clierrs

import (
	"errors"
)

var (
	ErrCallerNotFound                = errors.New("caller not found")
	ErrUserNotFound                  = errors.New("user not found")
	ErrUserAlreadyExist              = errors.New("user already exist")
	ErrInvalidValue                  = errors.New("value can't be setted to param")
	ErrInvalidAuthToken              = errors.New("invalid auth token")
	ErrTokenExpired                  = errors.New("token time to live is expired")
	ErrAuthTokenWasNotProvided       = errors.New("auth token required")
	ErrInvalidCredentials            = errors.New("incorrect password")
	ErrPermissionAdmin               = errors.New("not enough rights, you must to have admin role")
	ErrInvalidUserAttachedQuizFilter = errors.New("onlyCompleted and onlyNotCompleted could not be true at the same time")
	ErrInvalidUserCreatedQuizFilter  = errors.New("onlyPublished and onlyNotPublished could not be true at the same time")
	ErrInvalidUsersFilter            = errors.New("onlyCommon and onlyAdmins could not be true at the same time")
	ErrUpdateForbidden               = errors.New("you can't change this param")
	ErrInvalidUUID                   = errors.New("invalid user UUID")
)
