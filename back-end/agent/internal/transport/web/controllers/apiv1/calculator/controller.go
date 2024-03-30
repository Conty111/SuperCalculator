package calculator

import (
	"encoding/json"
	"github.com/Conty111/SuperCalculator/back-end/agent/internal/services"
	"github.com/Conty111/SuperCalculator/back-end/agent/internal/transport/web/controllers/apiv1"
	"github.com/Conty111/SuperCalculator/back-end/agent/internal/transport/web/helpers"
	"github.com/Conty111/SuperCalculator/back-end/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
)

var (
	_ apiv1.Controller = (*Controller)(nil)
)

// Controller is a controller implementation for status checks
type Controller struct {
	apiv1.BaseController
	Service      *services.CalculatorService
	RelativePath string
}

// NewController creates new status controller instance
func NewController(svc *services.CalculatorService) *Controller {
	return &Controller{
		Service:      svc,
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
	var body models.DurationSettings
	data, _ := io.ReadAll(ctx.Request.Body)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error().Err(err).Msg("incorrect closing request body")
		}
	}(ctx.Request.Body)
	err := json.Unmarshal(data, &body)
	if err != nil {
		helpers.WriteErrResponse(ctx, err)
		return
	}
	ctrl.Service.SetOperationDuration(&body)
	ctx.JSON(http.StatusOK, &Response{
		Status:  http.StatusText(http.StatusOK),
		Message: "time duration successfully updated",
	})
}

// DefineRoutes adds controller routes to the router
func (ctrl *Controller) DefineRoutes(r gin.IRouter) {
	r.PUT("/", ctrl.SetTime)
}
