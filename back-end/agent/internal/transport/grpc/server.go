package grpc

import (
	"context"
	"github.com/Conty111/SuperCalculator/back-end/agent/internal/services"
	"github.com/Conty111/SuperCalculator/back-end/models"
	pb "github.com/Conty111/SuperCalculator/back-end/proto"
)

type Server struct {
	pb.AgentGRPCServer
	Monitor    *services.Monitor
	Calculator *services.CalculatorService
}

func NewGRPCService(monitor *services.Monitor, calculator *services.CalculatorService) *Server {
	return &Server{
		Monitor:    monitor,
		Calculator: calculator,
	}
}

func (s *Server) GetInfo(ctx context.Context, req *pb.AgentInfoRequest) (*pb.AgentInfoResponse, error) {
	info := s.Monitor.GetInfo()
	settings := s.Calculator.GetSettings()
	return &pb.AgentInfoResponse{
		AgentID:        info.AgentID,
		Name:           info.Name,
		LastTaskID:     uint32(info.LastTaskID),
		CompletedTasks: uint32(info.CompletedTasks),
		Settings: &pb.AgentSettings{
			DivisionDuration: float32(settings.DivisionDuration),
			SubtractDuration: float32(settings.SubtractDuration),
			AddDuration:      float32(settings.AddDuration),
			MultiplyDuration: float32(settings.MultiplyDuration),
		},
	}, nil
}

func (s *Server) SetSettings(ctx context.Context, req *pb.SetAgentSettingsRequest) (*pb.SetAgentSettingsResponse, error) {
	s.Calculator.SetOperationDuration(&models.DurationSettings{
		DivisionDuration: float64(req.Settings.DivisionDuration),
		AddDuration:      float64(req.Settings.AddDuration),
		SubtractDuration: float64(req.Settings.SubtractDuration),
		MultiplyDuration: float64(req.Settings.MultiplyDuration),
	})
	return nil, nil
}
