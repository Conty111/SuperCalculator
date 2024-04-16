package grpc

import (
	"github.com/Conty111/SuperCalculator/back-end/models"
	pb "github.com/Conty111/SuperCalculator/back-end/proto"
)

type AgentGRPCClient struct {
	pb.AgentGRPCClient
}

func NewAgentGRPCClient() *AgentGRPCClient {
	return &AgentGRPCClient{}
}

func (c *AgentGRPCClient) GetAgentInfo(agent models.AgentConfig) (models.AgentInfo, error) {
	return models.AgentInfo{}, nil
}

func (c *AgentGRPCClient) SetAgentSettings(settings models.Settings, agent models.AgentConfig) error {
	//c.AgentGRPCClient.SetSettings()
	return nil
}
