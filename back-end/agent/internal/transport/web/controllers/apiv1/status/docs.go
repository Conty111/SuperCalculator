package status

// ResponseDoc is a response declaration for documentation purposes
type ResponseDoc struct {
	Data struct {
		Attributes Response `json:"attributes"`
	} `json:"data"`
}
