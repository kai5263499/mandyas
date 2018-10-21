package server

import (
	"bufio"
	"fmt"
	"io"
	"net"

	"github.com/kai5263499/mandyas/domain"
	pb "github.com/kai5263499/mandyas/generated"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Server grpc server
type Server struct {
	Conf          *domain.Config
	ServiceServer pb.MandyasServiceServer
	lis           *net.Listener
	grpcServer    *grpc.Server

	StdOutPipe io.Reader
	StdInPipe  io.WriteCloser

	cmdOutputChan chan string
	cmdInputChan  chan string
}

func New(config *domain.Config, stdOutPipe io.Reader, stdInPipe io.WriteCloser) (*Server, error) {
	cmdInputChan := make(chan string)

	serviceServer := MandyasService{
		CmdInputChan: cmdInputChan,
	}

	return &Server{
		Conf:          config,
		ServiceServer: serviceServer,

		StdOutPipe: stdOutPipe,
		StdInPipe:  stdInPipe,

		cmdOutputChan: make(chan string),
		cmdInputChan:  cmdInputChan,
	}, nil
}

func (s *Server) configureCmd() error {
	go func() {
		io.WriteString(s.StdInPipe, "values written to stdin are passed to cmd's standard input")
	}()

	return nil
}

func (s *Server) readCmdOutput() {
	rd := bufio.NewReader(s.StdOutPipe)

	for {
		str, err := rd.ReadString('\n')
		if err != nil {
			log.Fatal("Read Error:", err)
			return
		}
		trimmedStr := str[:len(str)-1]
		s.cmdOutputChan <- trimmedStr
	}
}

// Start the grpc service
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

	// log.Infof("starting read command output")
	// go s.readCmdOutput()

	log.Infof("now serving")
	go s.grpcServer.Serve(lis)

	return nil
}
