package enums

import "errors"

type TimeFormat string

const RFC3339 TimeFormat = "2006-01-02T15:04:05Z07:00"
const RFC3339Nano TimeFormat = "2006-01-02T15:04:05.999999999Z07:00"

var (
	ReadRepoErr   = errors.New("failed to read file or convert type")
	DeleteRepoErr = errors.New("failed to delete file or convert type")
)
