package user

import "github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/transport/web/helpers"

// ResponseDoc is a response declaration for documentation purposes
type ResponseDoc struct {
	Data struct {
		Attributes MsgResponse `json:"attributes"`
	} `json:"data"`
}

// MsgResponseDoc is a response declaration for documentation purposes
type MsgResponseDoc struct {
	Data struct {
		Attributes helpers.MsgResponse `json:"attributes"`
	} `json:"data"`
}
