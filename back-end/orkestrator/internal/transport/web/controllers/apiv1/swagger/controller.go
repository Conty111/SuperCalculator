package swagger

import (
	"github.com/Conty111/SuperCalculator/back-end/agent/internal/transport/web/controllers/apiv1"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

	//nolint: golint //reason: blank import because of swagger docs init
	_ "github.com/Conty111/SuperCalculator/back-end/agent/api/web"
)

var (
	_ apiv1.Controller = (*Controller)(nil)
)

// Controller implements controller for swagger
type Controller struct {
	apiv1.BaseController
	RelativePath string
}

// NewController create new instance for swagger controller
func NewController() *Controller {
	return &Controller{
		RelativePath: "/swagger",
	}
}

// GetRelativePath returns relative path of the controller's router group
func (ctrl *Controller) GetRelativePath() string {
	return ctrl.RelativePath
}

// DefineRoutes adds swagger controller routes to the router
func (ctrl *Controller) DefineRoutes(r gin.IRouter) {
	r.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
