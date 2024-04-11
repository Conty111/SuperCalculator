package helpers

import (
	"errors"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/clierrs"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AgentResponse struct {
	Body   map[string]interface{}
	Status int
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
	case errors.Is(err, clierrs.ErrTaskAlreadyCreated):
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
