package server

import (
	"context"

	pb "github.com/kai5263499/mandyas/generated"
)

// MandyasService represents a GRPC server
type MandyasService struct{}

// GetServerOutput starts a stream of server outputs to the given client
func (s MandyasService) GetServerOutput(_ *pb.GetServerOutputRequest, stream pb.MandyasService_GetServerOutputServer) error {
	// stream.Send(nil)
	return nil
}

// SendCommand sends a command to the underlying service
func (s MandyasService) SendCommand(ctx context.Context, command *pb.ServerCommandRequest) (*pb.ServerCommandResponse, error) {
	return nil, nil
}

// GetStatus returns underlying service status
func (s MandyasService) GetStatus(ctx context.Context, _ *pb.GetStatusRequest) (*pb.ServerStatus, error) {
	return nil, nil
}
