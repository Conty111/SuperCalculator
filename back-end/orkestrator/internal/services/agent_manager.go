package services

import (
	"github.com/Conty111/SuperCalculator/back-end/models"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/interfaces"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/transport/web/helpers"
)

type AgentManager struct {
	Agents []models.AgentConfig
	Client interfaces.AgentAPIClient
}

func NewAgentManager(
	agents []models.AgentConfig,
	client interfaces.AgentAPIClient) *AgentManager {

	return &AgentManager{
		Agents: agents,
		Client: client,
	}
}

func (s *AgentManager) SetSettings(settings *models.Settings) []*helpers.AgentResponse {
	//TODO implement me
	panic("implement me")
}

func (s *AgentManager) GetWorkersInfo() []*helpers.AgentResponse {
	panic("implement me")
}

func (s *AgentManager) GetAgentInfo(agent models.AgentConfig) (map[string]interface{}, int, error) {
	panic("implement me")
}
