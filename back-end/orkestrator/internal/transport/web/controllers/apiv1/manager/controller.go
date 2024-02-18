package manager

import (
	"errors"
	"github.com/Conty111/SuperCalculator/back-end/models"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/clierrs"
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
	Service      interfaces.Service
	RelativePath string
}

// NewController creates new status controller instance
func NewController(svc interfaces.Service) *Controller {
	return &Controller{
		Service:      svc,
		RelativePath: "/manager",
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
// @Router /api/v1/manager/settings [put]
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

// GetTasks godoc
// @Summary Returns All Tasks
// @Description set time duration
// @ID get-tasks
// @Accept json
// @Produce json
// @Success 200 {object} TasksListResponse
// @Router /api/v1/manager/tasks [get]
func (ctrl *Controller) GetTasks(ctx *gin.Context) {
	tasks, err := ctrl.Service.GetAllTasks()
	if err != nil {
		helpers.WriteErrResponse(ctx, err)
		return
	}
	results := make([]*Task, len(tasks))
	for i, t := range tasks {
		var res Task
		res.ID = t.ID
		res.Expression = t.Expression
		res.IsExecuted = t.IsExecuted
		res.CreatedAt = t.CreatedAt
		res.ExecutedAt = t.UpdatedAt
		res.Error = t.Error
		res.Value = t.Value
		results[i] = &res
	}
	ctx.JSON(http.StatusOK, &TasksListResponse{
		Status: http.StatusText(http.StatusOK),
		Tasks:  results,
	})
}

// HandleTask godoc
// @Summary HandleTask
// @Description saves task and sends it to execute
// @ID handle-task
// @Accept json
// @Produce json
// @Success 200 {object} ResponseDoc
// @Router /api/v1/manager [post]
func (ctrl *Controller) HandleTask(ctx *gin.Context) {
	var body models.Task
	if err := ctx.ShouldBind(&body); err != nil {
		helpers.WriteErrResponse(ctx, err)
		return
	}
	task, err := ctrl.Service.CreateTask(body.Expression)
	if err != nil {
		if errors.Is(err, clierrs.ErrTaskAlreadyCreated) {
			var resp struct {
				Response
				models.TasksModel
			}
			resp.Response = Response{
				Status:  http.StatusText(http.StatusOK),
				Message: "task already existed",
			}
			resp.TasksModel = *task
			ctx.JSON(http.StatusOK, &resp)
			return
		}
		helpers.WriteErrResponse(ctx, err)
		return
	}
	var resp struct {
		Response
		TaskID uint `json:"task_id"`
	}
	resp.Response = Response{
		Status:  http.StatusText(http.StatusOK),
		Message: "task successfully created",
	}
	resp.TaskID = task.ID
	ctx.JSON(http.StatusOK, &resp)
}

// GetWorkersInfo godoc
// @Summary Get Workers Info
// @Description returns info about all workers
// @ID workers-info
// @Accept json
// @Produce json
// @Success 200 {object} WorkersInfoResponse
// @Router /api/v1/manager/workers [get]
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
	r.POST("/", ctrl.HandleTask)
	r.GET("/workers", ctrl.GetWorkersInfo)
	r.GET("/tasks", ctrl.GetTasks)
}
