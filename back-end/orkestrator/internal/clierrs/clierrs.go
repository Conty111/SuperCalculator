package clierrs

import "errors"

var (
	AuthTokenWasNotProvided = errors.New("required jwt token")
	IvalidAuthToken         = errors.New("token is invalid")
	PermissionDenied        = errors.New("permission denied")
	FileAlreadyExist        = errors.New("file already exist")
	FileNotFound            = errors.New("file not found")
	FileQueryRequired       = errors.New("required file name or file id in query params")
)
