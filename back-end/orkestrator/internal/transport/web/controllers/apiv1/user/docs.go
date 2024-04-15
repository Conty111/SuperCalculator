package user

import "gitlab.sch.ocrv.com.rzd/blockchain/trainees/quiz-svc/internal/web/helpers"

// ResponseDoc is a response declaration for documentation purposes
type ResponseDoc struct {
	Data struct {
		Attributes Response `json:"attributes"`
	} `json:"data"`
}

// MsgResponseDoc is a response declaration for documentation purposes
type MsgResponseDoc struct {
	Data struct {
		Attributes helpers.MsgResponse `json:"attributes"`
	} `json:"data"`
}
