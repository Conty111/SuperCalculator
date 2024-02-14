package router

import (
	"github.com/gin-gonic/gin"
)

// NewRouter create new gin router instance
func NewRouter() *gin.Engine {
	r := gin.New()

	return r
}
