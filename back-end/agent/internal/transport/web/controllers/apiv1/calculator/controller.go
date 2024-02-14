package calculator

import (
	"github.com/Conty111/SuperCalculator/back-end/agent/internal/services"
	"github.com/Conty111/SuperCalculator/back-end/agent/internal/transport/web/controllers/apiv1"
	"github.com/Conty111/SuperCalculator/back-end/agent/internal/transport/web/helpers"
	"github.com/Conty111/SuperCalculator/back-end/agent/internal/transport/web/render"
	"github.com/gin-gonic/gin"

	"net/http"
)

var (
	_ apiv1.Controller = (*Controller)(nil)
)

// Controller is a controller implementation for status checks
type Controller struct {
	apiv1.BaseController
	Service      *services.ExpressionService
	RelativePath string
}

// NewController creates new status controller instance
func NewController(svc *services.ExpressionService) *Controller {
	return &Controller{
		RelativePath: "/calculator",
	}
}

// GetRelativePath returns relative path of the controller's router group
func (ctrl *Controller) GetRelativePath() string {
	return ctrl.RelativePath
}

// SetTime godoc
// @Summary Set Time To Execute Operation
// @Description set time duration
// @ID set-time
// @Accept json
// @Produce json
// @Success 200 {object} ResponseDoc
// @Router /api/v1/calculator [put]
func (ctrl *Controller) SetTime(ctx *gin.Context) {
	var body RequestBody
	if err := ctx.ShouldBind(&body); err != nil {
		helpers.WriteErrResponse(ctx, err)
		return
	}
	err := ctrl.Service.SetOperationDuration(body.Operation, body.Duration)
	if err != nil {
		helpers.WriteErrResponse(ctx, err)
		return
	}
	render.JSONAPIPayload(ctx, http.StatusOK, &Response{
		Status:  http.StatusText(http.StatusOK),
		Message: "time successfully updated",
	})
}

// DefineRoutes adds controller routes to the router
func (ctrl *Controller) DefineRoutes(r gin.IRouter) {
	r.PUT("/", ctrl.SetTime)
}
