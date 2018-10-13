package domain

import (
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/kai5263499/mandyas/generated"
)

type Server struct {
	Conf          *Config
	ServiceServer pb.MandyasServiceServer
	lis           *net.Listener
	grpcServer    *grpc.Server
}

func (s *Server) Start() error {
	listenAddress := fmt.Sprintf("localhost:%d", s.Conf.GrpcPort)

	lis, err := net.Listen("tcp", listenAddress)

	if err != nil {
		log.Errorf("unable to create network listener err=%#v", err)
		return err
	}

	s.lis = &lis

	log.Infof("listening on %s", listenAddress)

	if !s.Conf.UseTLS {
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
			return err
		}
		s.grpcServer = grpc.NewServer()
	} else {
		certFile := s.Conf.SSLCertFile
		keyFile := s.Conf.SSLKeyFile
		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
			return err
		}
		s.grpcServer = grpc.NewServer(grpc.Creds(creds))
	}

	pb.RegisterMandyasServiceServer(s.grpcServer, s.ServiceServer)

	log.Infof("now serving")

	return s.grpcServer.Serve(lis)
}
