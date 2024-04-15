package user

import (
	"gitlab.sch.ocrv.com.rzd/blockchain/trainees/quiz-svc/internal/enums"
	"gitlab.sch.ocrv.com.rzd/blockchain/trainees/quiz-svc/internal/web/helpers"
	"gitlab.sch.ocrv.com.rzd/blockchain/trainees/quiz-svc/internal/web/serializers"
)

type MsgResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	UserID  uint   `json:"user_id"`
}

// Response is a declaration for a user response
type Response struct {
	Status string `json:"status"`
	serializers.Info
}

// SuperUserResponse is a declaration for response to a superuser
type SuperUserResponse struct {
	Status string `json:"status"`
	Info   serializers.VerboseInfo
}

// AuthResponse is a declaration for an auth response
type AuthResponse struct {
	Status string      `json:"status"`
	Role   enums.Role  `json:"role"`
	Theme  enums.Theme `json:"theme"`
	Token  string      `json:"token"`
}

// AllUsersResponse is a declaration for an all users response
type AllUsersResponse struct {
	Status     string                 `json:"status"`
	Pagination helpers.PaginationInfo `json:"pagination"`
	Data       []*serializers.Info    `json:"users"`
}

/////////////////////////////////// Responses for the embedded resources /////////////////////////////////////////

type AttachedQuizzesResponse struct {
	Status     string                               `json:"status"`
	Pagination helpers.PaginationInfo               `json:"pagination"`
	Data       []*serializers.ShortUserAttachedQuiz `json:"attachedQuizzes"`
}

type CreatedQuizzesResponse struct {
	Status     string                   `json:"status"`
	Pagination helpers.PaginationInfo   `json:"pagination"`
	Data       []*serializers.ShortQuiz `json:"createdQuizzes"`
}
