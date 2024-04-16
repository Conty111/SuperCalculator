package grpc

import (
	"context"
	"fmt"
	"github.com/Conty111/SuperCalculator/back-end/models"
	pb "github.com/Conty111/SuperCalculator/back-end/proto"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AgentGRPCClient struct{}

func NewAgentGRPCClient() *AgentGRPCClient {
	return &AgentGRPCClient{}
}

func (c *AgentGRPCClient) GetAgentInfo(agent *models.AgentConfig) (*models.AgentInfo, error) {
	conn, err := c.connect(agent)
	defer conn.Close()
	grpcClient := pb.NewAgentGRPCClient(conn)
	info, err := grpcClient.GetInfo(context.TODO(), &pb.AgentInfoRequest{})
	if err != nil {
		return nil, err
	}
	return &models.AgentInfo{
		Name:           info.Name,
		AgentID:        info.AgentID,
		CompletedTasks: uint(info.CompletedTasks),
		LastTaskID:     uint(info.LastTaskID),
	}, nil
}

func (c *AgentGRPCClient) SetAgentSettings(settings *models.Settings, agent *models.AgentConfig) error {
	conn, err := c.connect(agent)
	defer conn.Close()
	grpcClient := pb.NewAgentGRPCClient(conn)
	_, err = grpcClient.SetSettings(
		context.TODO(),
		&pb.SetAgentSettingsRequest{
			Settings: &pb.AgentSettings{
				DivisionDuration: float32(settings.DivisionDuration),
				AddDuration:      float32(settings.AddDuration),
				SubtractDuration: float32(settings.SubtractDuration),
				MultiplyDuration: float32(settings.MultiplyDuration),
			},
		})
	if err != nil {
		return err
	}
	return nil
}

func (c *AgentGRPCClient) connect(agent *models.AgentConfig) (*grpc.ClientConn, error) {
	addr := fmt.Sprintf("%s:%d", agent.Address, agent.GrpcPort) // используем адрес сервера
	// установим соединение
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error().Err(err).Msg("could not connect to grpc server")
		return nil, err
	}
	return conn, nil
}
