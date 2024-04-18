package services

import (
	"github.com/Conty111/SuperCalculator/back-end/models"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/interfaces"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/transport/web/helpers"
	"github.com/rs/zerolog/log"
	"go/types"
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

func (s *AgentManager) SetSettings(settings *models.Settings) []*helpers.AgentResponse[types.Nil] {
	//wg := sync.WaitGroup{}
	responses := make([]*helpers.AgentResponse[types.Nil], len(s.Agents))

	for i, agent := range s.Agents {
		//wg.Add(1)
		//agent := agent
		//i := i
		//go func() {
		//	defer wg.Done()
		var res helpers.AgentResponse[types.Nil]

		err := s.Client.SetAgentSettings(settings, &agent)
		if err != nil {
			log.Error().Err(err).
				Str("agent", agent.Name).
				Str("agentAddress", agent.Address).
				Msg("error while setting settings to the agent")
			res.Error = err.Error()
		}
		responses[i] = &res
		//}()
	}
	//wg.Wait()
	return responses
}

func (s *AgentManager) GetWorkersInfo() []*helpers.AgentResponse[*models.AgentInfo] {
	//wg := sync.WaitGroup{}
	responses := make([]*helpers.AgentResponse[*models.AgentInfo], len(s.Agents))
	for i, agent := range s.Agents {
		//wg.Add(1)
		//agent := agent
		//i := i
		//go func() {
		//	defer wg.Done()
		var res helpers.AgentResponse[*models.AgentInfo]
		info, err := s.Client.GetAgentInfo(&agent)
		if err != nil {
			log.Error().Err(err).
				Str("agent", agent.Name).
				Str("agentAddress", agent.Address).
				Msg("error while getting info from agent")
			res.Error = err.Error()
		}
		res.Body = info
		responses[i] = &res
		//	}()
	}
	//wg.Wait()
	return responses
}
