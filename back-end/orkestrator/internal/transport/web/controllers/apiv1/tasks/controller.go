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

func (ctrl *Controller) GetTasks(ctx *gin.Context) {
	tasks, err := ctrl.Service.GetAllTasks()
	if err != nil {
		helpers.WriteErrResponse(ctx, err)
		return
	}
	results := make([]*helpers.TaskResponse, len(tasks))
	for i, t := range tasks {
		var res helpers.TaskResponse
		res.ID = t.ID
		res.Expression = t.Expression
		res.IsExecuted = t.IsExecuted
		res.CreatedAt = t.CreatedAt
		res.ExecutedAt = t.UpdatedAt
		res.Error = t.Error
		res.Value = t.Value
		res.UserID = t.User.ID
		results[i] = &res
	}
	ctx.JSON(http.StatusOK, &TasksListResponse{
		Status: http.StatusText(http.StatusOK),
		Tasks:  results,
	})
}

func (ctrl *Controller) HandleTask(ctx *gin.Context) {
	var body models.Task
	if err := ctx.ShouldBind(&body); err != nil {
		helpers.WriteErrResponse(ctx, err)
		return
	}
	callerID := ctx.GetUint("callerID")
	task, err := ctrl.Service.CreateTask(body.Expression, callerID)
	if err != nil {
		if errors.Is(err, clierrs.ErrTaskAlreadyCreated) {
			ctx.JSON(http.StatusOK, &Response{
				Status:  http.StatusText(http.StatusOK),
				Message: "task already existed",
			})
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
	r.POST("", ctrl.HandleTask)
	r.GET("", ctrl.GetTasks)
}
