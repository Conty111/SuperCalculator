package agents_manager

import (
	"github.com/Conty111/SuperCalculator/back-end/models"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/interfaces"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/transport/web/controllers/apiv1"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/transport/web/helpers"
	"github.com/gin-gonic/gin"
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

// SetSettings godoc
// @Summary Set Calculation Settings
// @Description set time duration to execute and timeout to workers
// @ID set-settings
// @Accept json
// @Produce json
// @Success 200 {object} ResponseDoc
// @Router /api/v1/tasks/settings [put]
func (ctrl *Controller) SetSettings(ctx *gin.Context) {
	var body models.Settings
	if err := ctx.ShouldBind(&body); err != nil {
		helpers.WriteErrResponse(ctx, err)
		return
	}
	responses, statuses := ctrl.Service.SetSettings(&body)
	ctx.JSON(http.StatusOK, &WorkersListResponse{
		Status:    http.StatusText(http.StatusOK),
		Responses: serializeWorkersResponse(responses, statuses),
	})
}

// GetWorkersInfo godoc
// @Summary Get Workers Info
// @Description returns info about all workers
// @ID workers-info
// @Accept json
// @Produce json
// @Success 200 {object} WorkersInfoResponse
// @Router /api/v1/tasks/workers [get]
func (ctrl *Controller) GetWorkersInfo(ctx *gin.Context) {
	responses, statuses := ctrl.Service.GetWorkersInfo()
	ctx.JSON(http.StatusOK, &WorkersListResponse{
		Status:    http.StatusText(http.StatusOK),
		Responses: serializeWorkersResponse(responses, statuses),
	})
}

func serializeWorkersResponse(responses []map[string]interface{}, statuses []int) []WorkerResponse {
	workersResponses := make([]WorkerResponse, len(responses))
	for i, resp := range responses {
		workersResponses[i] = WorkerResponse{
			Status:   http.StatusText(statuses[i]),
			Response: resp,
		}
	}
	return workersResponses
}

// DefineRoutes adds controller routes to the router
func (ctrl *Controller) DefineRoutes(r gin.IRouter) {
	r.PUT("/settings", ctrl.SetSettings)
	r.GET("/info", ctrl.GetWorkersInfo)
}
