package grpcsrv

import (
	"context"
	"errors"
	"net"

	"github.com/nsmak/otus_hw/hw12_13_14_15_calendar/internal/app"
	"google.golang.org/grpc"
)

//go:generate protoc ./proto/EventService.proto --go_out=. --go-grpc_out=.

type ServerError struct {
	Message string `json:"message"`
	Err     error  `json:"err,omitempty"`
}

func (e *ServerError) Error() string {
	if e.Err != nil {
		e.Message = e.Message + " --> " + e.Err.Error()
	}
	return "[grpc] " + e.Message
}
func (e *ServerError) Unwrap() error {
	return e.Err
}

type Server struct {
	Address string
	server  *grpc.Server
	api     *API
	log     app.Logger
}

func NewServer(api *API, host, port string, logger app.Logger) *Server {
	return &Server{
		Address: net.JoinHostPort(host, port),
		api:     api,
		log:     logger,
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.server = grpc.NewServer(grpc.UnaryInterceptor(s.loggingInterceptor))
	RegisterEventServiceServer(s.server, s.api)
	lis, err := net.Listen("tcp", s.Address)
	if err != nil {
		return &ServerError{Message: "start listen error", Err: err}
	}

	if err = s.server.Serve(lis); err != nil {
		return &ServerError{Message: "start server error", Err: err}
	}

	<-ctx.Done()
	return nil
}

func (s *Server) Stop() error {
	if s.server == nil {
		return errors.New("grpc server is nil")
	}

	s.server.GracefulStop()
	return nil
}
