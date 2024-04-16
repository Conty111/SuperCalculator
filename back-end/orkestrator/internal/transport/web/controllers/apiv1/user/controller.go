package user

import (
	"github.com/Conty111/SuperCalculator/back-end/models"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/transport/web/controllers/apiv1"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/transport/web/helpers"
	"github.com/cristalhq/jwt/v5"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var (
	_ apiv1.Controller = (*Controller)(nil)
)

type Service interface {
	GetUserByID(userID, callerID uint) (*models.User, error)
	UpdateUserParamByID(userID uint, param string, value interface{}, callerID uint) error
	DeleteUserByID(userID uint, callerID uint) error
	GetAllUsers(callerID uint) ([]*models.User, error)

	CreateUser(user *models.User) (*jwt.Token, error)
	Login(email, password string) (*models.User, *jwt.Token, error)
	//GetMeByID(callerID uint) (*models.User, error)
	//UpdateMeByID(callerID uint, param, value string) error
}

// Controller is a controller implementation for user endpoints
type Controller struct {
	apiv1.BaseController
	Service Service
}

// NewController creates new user controller instance
func NewController(service Service) *Controller {
	return &Controller{
		BaseController: apiv1.BaseController{
			RelativePath: "/users",
		},
		Service: service,
	}
}

// GetRelativePath returns relative path of the controller's router group
func (ctrl *Controller) GetRelativePath() string {
	return ctrl.RelativePath
}

// GetUser endpoint
func (ctrl *Controller) GetUser(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("userID"))
	if err != nil {
		helpers.WriteErrResponse(ctx, err)
		return
	}
	callerID := ctx.GetUint("callerID")

	user, err := ctrl.Service.GetUserByID(uint(userID), callerID)
	if err != nil {
		helpers.WriteErrResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, &UserResponse{
		Status: http.StatusText(http.StatusOK),
		UserInfo: UserInfo{
			Role:     user.Role,
			Username: user.Username,
			Email:    user.Email,
		},
	})
}

// CreateUser endpoint
func (ctrl *Controller) CreateUser(ctx *gin.Context) {
	var user models.User

	if err := ctx.ShouldBind(&user); err != nil {
		helpers.WriteErrResponse(ctx, err)
		return
	}

	token, err := ctrl.Service.CreateUser(&user)
	if err != nil {
		helpers.WriteErrResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, &AuthResponse{
		Status: http.StatusText(http.StatusOK),
		Token:  token.String(),
		UserInfo: UserInfo{
			Role:     user.Role,
			Username: user.Username,
			Email:    user.Email,
		},
	})
}

// Upd is body of request to update endpoint
type Upd struct {
	Param string      `json:"param"`
	Value interface{} `json:"value"`
}

// UpdateUser endpoint
func (ctrl *Controller) UpdateUser(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("userID"))
	if err != nil {
		helpers.WriteErrResponse(ctx, err)
		return
	}
	callerID := ctx.GetUint("callerID")

	var upd Upd
	err = ctx.ShouldBind(&upd)
	if err != nil {
		helpers.WriteErrResponse(ctx, err)
		return
	}

	err = ctrl.Service.UpdateUserParamByID(uint(userID), upd.Param, upd.Value, callerID)
	if err != nil {
		helpers.WriteErrResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, &helpers.MsgResponse{
		Status:  http.StatusText(http.StatusOK),
		Message: "user successfully updated",
	})
}

// DeleteUser endpoint
func (ctrl *Controller) DeleteUser(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("userID"))
	if err != nil {
		helpers.WriteErrResponse(ctx, err)
		return
	}
	callerID := ctx.GetUint("callerID")

	err = ctrl.Service.DeleteUserByID(uint(userID), callerID)
	if err != nil {
		helpers.WriteErrResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, &helpers.MsgResponse{
		Status:  http.StatusText(http.StatusOK),
		Message: "user successfully deleted",
	})
}

// GetAllUsers endpoint
func (ctrl *Controller) GetAllUsers(ctx *gin.Context) {
	callerID := ctx.GetUint("callerID")

	users, err := ctrl.Service.GetAllUsers(callerID)
	if err != nil {
		helpers.WriteErrResponse(ctx, err)
		return
	}
	usersInfoList := make([]UserInfo, len(users))
	for i, u := range users {
		usersInfoList[i] = UserInfo{
			Username: u.Username,
			Email:    u.Email,
			Role:     u.Role,
		}
	}

	ctx.JSON(http.StatusOK, &UsersListResponse{
		Status: http.StatusText(http.StatusOK),
		Users:  usersInfoList,
	})
}

// AuthRequestBody is a declaration for an auth request body
type AuthRequestBody struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Login endpoint
func (ctrl *Controller) Login(ctx *gin.Context) {
	var body AuthRequestBody
	if err := ctx.ShouldBind(&body); err != nil {
		helpers.WriteErrResponse(ctx, err)
		return
	}
	user, token, err := ctrl.Service.Login(body.Email, body.Password)
	if err != nil {
		helpers.WriteErrResponse(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, &AuthResponse{
		Status: http.StatusText(http.StatusOK),
		Token:  token.String(),
		UserInfo: UserInfo{
			Role:     user.Role,
			Username: user.Username,
			Email:    user.Username,
		},
	})
}

// GetMe godoc
// @Tags users,me
// @Summary Get User
// @Description get me (using id from token)
// @ID get-me
// @Accept json
// @Produce json
// @Success 200 {object} Response
// @Router /api/v1/users/me [get]
//func (ctrl *Controller) GetMe(ctx *gin.Context) {
//	callerID, err := uuid.FromString(ctx.GetString("callerID"))
//	if err != nil {
//		helpers.WriteErrResponse(ctx, err)
//		return
//	}
//
//	user, err := ctrl.Service.GetMeByID(callerID)
//	if err != nil {
//		helpers.WriteErrResponse(ctx, err)
//		return
//	}
//
//	ctx.JSON(http.StatusOK, &Response{
//		Status: http.StatusText(http.StatusOK),
//		Info: serializers.Info{
//			ID:       callerID,
//			FullName: user.FullName,
//			Email:    user.Email,
//			Role:     user.Role,
//			Theme:    user.Theme,
//		},
//	})
//}

// UpdateMe godoc
// @Tags users,me
// @Summary Update Me
// @Description update user param by id from token
// @ID update-me
// @Accept json
// @Produce json
// @Param upd     body helpers.Upd true "A structure consisting of the parameter being updated and its new value."
// @Success 200 {object} helpers.MsgResponse
// @Failure 400 {object} helpers.ErrResponse
// @Failure 403 {object} helpers.ErrResponse
// @Router /api/v1/users/me [patch]
//func (ctrl *Controller) UpdateMe(ctx *gin.Context) {
//	callerID, err := uuid.FromString(ctx.GetString("callerID"))
//	if err != nil {
//		helpers.WriteErrResponse(ctx, err)
//		return
//	}
//
//	var upd helpers.Upd
//	err = ctx.ShouldBind(&upd)
//	if err != nil {
//		helpers.WriteErrResponse(ctx, err)
//		return
//	}
//
//	err = ctrl.Service.UpdateMeByID(callerID, upd.Param, upd.Value)
//	if err != nil {
//		helpers.WriteErrResponse(ctx, err)
//		return
//	}
//
//	ctx.JSON(http.StatusOK, &helpers.MsgResponse{
//		Status:  http.StatusText(http.StatusOK),
//		Message: "user successfully updated",
//	})
//}

// DefineRoutes adds controller routes to the router
func (ctrl *Controller) DefineRoutes(r gin.IRouter) {
	// CRUD available only for superusers (except creation aka registration)
	r.POST("/create", ctrl.CreateUser)
	r.GET("/:userID", ctrl.GetUser)
	r.PATCH("/:userID", ctrl.UpdateUser)
	r.DELETE("/:userID", ctrl.DeleteUser)

	//// extra collections
	r.GET("", ctrl.GetAllUsers)

	// Available for everyone
	//r.GET("/me", ctrl.GetMe)
	//r.PATCH("/me", ctrl.UpdateMe)
	//r.GET("/me/attachedQuizzes", ctrl.GetMyAttachedQuizzes)
	//r.GET("/me/createdQuizzes", ctrl.GetMyCreatedQuizzes)
	r.POST("/login", ctrl.Login)
}
