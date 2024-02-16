package calculator

// Response is a declaration for a calculator response
type Response struct {
	Status  string `json:"attr,status"`
	Message string `json:"message"`
}
