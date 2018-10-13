package server

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/kai5263499/mandyas/generated"
)

// MandyasServer represents a GRPC server
type MandyasService struct{}

// GetServerOutput starts a stream of server outputs to the given client
func (s MandyasService) GetServerOutput(_ *empty.Empty, stream pb.MandyasService_GetServerOutputServer) error {
	// stream.Send(nil)
	return nil
}

// SendCommand sends a command to the underlying service
func (s MandyasService) SendCommand(ctx context.Context, command *pb.ServerCommand) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

// GetStatus returns underlying service status
func (s MandyasService) GetStatus(ctx context.Context, _ *empty.Empty) (*pb.ServerStatus, error) {
	return nil, nil
}
