package models

type DurationSettings struct {
	DivisionDuration float64 `json:"division_time"`
	AddDuration      float64 `json:"add_time"`
	MultiplyDuration float64 `json:"multiply_time"`
	SubtractDuration float64 `json:"subtract_time"`
}

type Settings struct {
	DurationSettings
	TimeoutResponse uint `json:"timeout_response"`
	TimeToRetry     uint `json:"time_retry"`
}
