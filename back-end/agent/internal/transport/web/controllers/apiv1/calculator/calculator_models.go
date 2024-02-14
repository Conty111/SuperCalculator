package calculator

import "time"

// Response is a declaration for a status response
type Response struct {
	Status  string `jsonapi:"attr,status"`
	Message string `json:"message"`
}

type RequestBody struct {
	Operation rune          `json:"operation"`
	Duration  time.Duration `json:"duration"`
}
