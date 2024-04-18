package agents_manager

import (
	"github.com/Conty111/SuperCalculator/back-end/models"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/interfaces"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/transport/web/controllers/apiv1"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/transport/web/helpers"
	"github.com/gin-gonic/gin"
	"go/types"
	"net/http"
)

var (
	_ apiv1.Controller = (*Controller)(nil)
)

// Controller is a controller implementation for status checks
type Controller struct {
	apiv1.BaseController
	Service      interfaces.AgentManager
	RelativePath string
}

// NewController creates new status controller instance
func NewController(svc interfaces.AgentManager) *Controller {
	return &Controller{
		Service:      svc,
		RelativePath: "/workers",
	}
}

// GetRelativePath returns relative path of the controller's router group
func (ctrl *Controller) GetRelativePath() string {
	return ctrl.RelativePath
}

func (ctrl *Controller) SetSettings(ctx *gin.Context) {
	var body models.Settings
	if err := ctx.ShouldBind(&body); err != nil {
		helpers.WriteErrResponse(ctx, err)
		return
	}
	responses := ctrl.Service.SetSettings(&body)
	ctx.JSON(http.StatusOK, &WorkersListResponse[types.Nil]{
		Responses: responses,
	})
}

func (ctrl *Controller) GetWorkersInfo(ctx *gin.Context) {
	responses := ctrl.Service.GetWorkersInfo()
	ctx.JSON(http.StatusOK, &WorkersListResponse[*models.AgentInfo]{
		Responses: responses,
	})
}

// DefineRoutes adds controller routes to the router
func (ctrl *Controller) DefineRoutes(r gin.IRouter) {
	r.PUT("/settings", ctrl.SetSettings)
	r.GET("/info", ctrl.GetWorkersInfo)
}
