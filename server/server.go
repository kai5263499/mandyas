package server

import (
	"context"
	"sync"

	pb "github.com/kai5263499/mandyas/generated"
	"github.com/satori/go.uuid"
)

var outputRecieverLock = &sync.Mutex{}

// MandyasService represents a GRPC server
type MandyasService struct {
	CmdInputChan              chan string
	registeredOutputRecievers map[uuid.UUID]pb.MandyasService_GetServerOutputServer
}

// SendCmdOutput sends command output to listeners
func (s MandyasService) SendCmdOutput(o string) int {
	sent := 0

	for i := range s.registeredOutputRecievers {
		so := &pb.ServerOutput{
			Id:      0,
			Content: []byte(o),
		}
		s.registeredOutputRecievers[i].Send(so)
		sent++
	}

	return sent
}

// GetServerOutput starts a stream of server outputs to the given client
func (s MandyasService) GetServerOutput(_ *pb.GetServerOutputRequest, stream pb.MandyasService_GetServerOutputServer) error {
	// stream.Send(nil)

	outputRecieverLock.Lock()
	id := uuid.NewV4()
	s.registeredOutputRecievers[id] = stream
	outputRecieverLock.Unlock()

	return nil
}

// SendCommand sends a command to the underlying service
func (s MandyasService) SendCommand(ctx context.Context, command *pb.ServerCommandRequest) (*pb.ServerCommandResponse, error) {
	s.CmdInputChan <- string(command.Command)
	return &pb.ServerCommandResponse{}, nil
}

// GetStatus returns underlying service status
func (s MandyasService) GetStatus(ctx context.Context, _ *pb.GetStatusRequest) (*pb.ServerStatus, error) {
	return &pb.ServerStatus{
		Status: pb.ServerStatus_RUNNING,
	}, nil
}
