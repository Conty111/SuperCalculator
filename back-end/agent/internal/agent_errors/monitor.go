package agent_errors

import "errors"

var (
	ErrTaskNotFound = errors.New("task with this id not found")
)
