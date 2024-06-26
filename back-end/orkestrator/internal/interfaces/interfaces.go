package interfaces

import (
	"context"
	"github.com/Conty111/SuperCalculator/back-end/models"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/transport/web/helpers"
	"github.com/cristalhq/jwt/v5"
	"go/types"
	"time"
)

type AgentManager interface {
	SetSettings(settings *models.Settings) []*helpers.AgentResponse[types.Nil]
	GetWorkersInfo() []*helpers.AgentResponse[*models.AgentInfo]
}

type AgentAPIClient interface {
	GetAgentInfo(agent *models.AgentConfig) (*models.AgentInfo, error)
	SetAgentSettings(settings *models.Settings, agent *models.AgentConfig) error
}

type TaskManager interface {
	GetAllTasks() ([]*models.TasksModel, error)
	CreateTask(expression string, callerID uint) (*models.TasksModel, error)
	DeleteTaskByID(taskID uint) error
	DeleteTaskByExpression(expression string) error
	SaveResult(res *models.Result) error
	Start(ctx context.Context)
}

//go:generate go run github.com/vektra/mockery/v3 --name UserProvider
type UserProvider interface {
	UserExists(userID uint) (bool, error)
	GetUserByID(userID uint) (*models.User, error)
}

//go:generate go run github.com/vektra/mockery/v3 --name UserManager
type UserManager interface {
	CreateUser(user *models.User) error
	GetUserByID(userID uint) (*models.User, error)
	GetUserByEmail(userEmail string) (*models.User, error)
	UserEmailExists(userEmail string) (bool, error)
	UserExists(userID uint) (bool, error)
	UpdateUser(user *models.User, param string, value interface{}) error
	DeleteUser(user *models.User) error
	GetLastID() (uint, error)
	GetAllUsers(callerID uint) ([]*models.User, error)
}

//go:generate go run github.com/vektra/mockery/v3 --name AuthManager
type AuthManager interface {
	BuildToken(userID uint) (*jwt.Token, error)
	GetTokenTTL() time.Duration
	HashString(text string) (string, error)
}
