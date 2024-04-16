package helpers

import (
	"errors"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/clierrs"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"net/http"
)

type AgentResponse struct {
	Body   map[string]interface{} `json:"body"`
	Status int                    `json:"status"`
}

type ErrResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

type MsgResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// WriteErrResponse writes an error response to the context
func WriteErrResponse(ctx *gin.Context, err error) {
	var status int
	switch {
	case errors.Is(err, clierrs.ErrTaskAlreadyCreated),
		errors.Is(err, clierrs.ErrInvalidValue),
		errors.Is(err, clierrs.ErrInvalidAuthToken),
		errors.Is(err, clierrs.ErrInvalidUUID),
		errors.Is(err, clierrs.ErrInvalidUserAttachedQuizFilter),
		errors.Is(err, clierrs.ErrInvalidUserCreatedQuizFilter),
		errors.Is(err, clierrs.ErrInvalidUsersFilter),
		errors.Is(err, clierrs.ErrUserAlreadyExist):

		status = http.StatusBadRequest
	case errors.Is(err, clierrs.ErrPermissionAdmin),
		errors.Is(err, clierrs.ErrAuthTokenWasNotProvided),
		errors.Is(err, clierrs.ErrTokenExpired),
		errors.Is(err, clierrs.ErrUpdateForbidden),
		errors.Is(err, clierrs.ErrInvalidCredentials):

		log.Error().Err(err).Msg("Attempt to get not allowed content")
		status = http.StatusForbidden

	case errors.Is(err, clierrs.ErrUserNotFound),
		errors.Is(err, clierrs.ErrCallerNotFound),
		errors.Is(err, clierrs.ErrInvalidCredentials):

		status = http.StatusNotFound
	default:
		status = http.StatusInternalServerError
	}

	if errors.As(err, &validator.ValidationErrors{}) {
		status = http.StatusBadRequest
	}

	ctx.AbortWithStatusJSON(status, &ErrResponse{
		Status: http.StatusText(status),
		Error:  err.Error(),
	})
}
