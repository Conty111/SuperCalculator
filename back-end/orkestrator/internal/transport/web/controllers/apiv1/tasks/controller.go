package tasks

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
	Service      interfaces.TaskManager
	RelativePath string
}

// NewController creates new status controller instance
func NewController(svc interfaces.TaskManager) *Controller {
	return &Controller{
		Service:      svc,
		RelativePath: "/tasks",
	}
}

// GetRelativePath returns relative path of the controller's router group
func (ctrl *Controller) GetRelativePath() string {
	return ctrl.RelativePath
}

// GetTasks godoc
// @Summary Returns All Tasks
// @Description set time duration
// @ID get-tasks
// @Accept json
// @Produce json
// @Success 200 {object} TasksListResponse
// @Router /api/v1/tasks/tasks [get]
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
// @Router /api/v1/tasks [post]
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

// DefineRoutes adds controller routes to the router
func (ctrl *Controller) DefineRoutes(r gin.IRouter) {
	r.POST("/execute", ctrl.HandleTask)
	r.GET("/tasks", ctrl.GetTasks)
}
