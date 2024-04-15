package clierrs

import "errors"

var (
	ErrTaskNotFound       = errors.New("task not found")
	ErrTaskAlreadyCreated = errors.New("task already created")
)
