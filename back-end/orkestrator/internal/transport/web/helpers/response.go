package helpers

import (
	"errors"
	"github.com/Conty111/SuperCalculator/back-end/agent/internal/clierrs"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ErrResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

type MsgResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	UserId  uint   `json:"userID"`
}

// WriteErrResponse writes an error response to the context
func WriteErrResponse(ctx *gin.Context, err error) {
	var status int
	switch {
	case errors.Is(err, clierrs.FileNotFound), errors.Is(err, clierrs.FileAlreadyExist):
		status = http.StatusNotFound
	case errors.Is(err, clierrs.PermissionDenied):
		status = http.StatusForbidden
	case errors.Is(err, clierrs.AuthTokenWasNotProvided):
		status = http.StatusBadRequest
	default:
		status = http.StatusInternalServerError
	}

	ctx.Error(err)
	ctx.AbortWithStatusJSON(status, &ErrResponse{
		Status: http.StatusText(status),
		Error:  err.Error(),
	})
}
