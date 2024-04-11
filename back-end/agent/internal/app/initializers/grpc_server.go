package initializers

import (
	"github.com/Conty111/SuperCalculator/back-end/agent/internal/app/dependencies"
	g "github.com/Conty111/SuperCalculator/back-end/agent/internal/transport/grpc"
	pb "github.com/Conty111/SuperCalculator/back-end/proto"
	"google.golang.org/grpc"
)

// InitializeGRPCServer create new grpc Server instance
func InitializeGRPCServer(container *dependencies.Container) *grpc.Server {
	server := grpc.NewServer()
	svc := g.NewGRPCService(container.Monitor, container.Calculator)
	pb.RegisterAgentGRPCServer(server, svc)
	return server
}
