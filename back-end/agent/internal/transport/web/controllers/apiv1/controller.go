package apiv1

import (
	"github.com/gin-gonic/gin"
)

// Controller is an interface for HTTP controllers
type Controller interface {
	DefineRoutes(gin.IRouter)
	// GetRelativePath returns relative path of the controller's router group
	GetRelativePath() string
}
