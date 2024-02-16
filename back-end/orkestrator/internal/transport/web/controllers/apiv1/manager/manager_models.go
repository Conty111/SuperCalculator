package calculator

// Response is a declaration for a status response
type Response struct {
	Status  string `jsonapi:"attr,status"`
	Message string `json:"message"`
}

type RequestBody struct {
	Operation string  `json:"operation"`
	Duration  float64 `json:"duration"`
}
