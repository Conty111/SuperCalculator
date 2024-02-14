package helpers

import (
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
	switch err {
	default:
		status = http.StatusInternalServerError
	}

	ctx.Error(err)
	ctx.AbortWithStatusJSON(status, &ErrResponse{
		Status: http.StatusText(status),
		Error:  err.Error(),
	})
}
