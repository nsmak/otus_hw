package grpcsrv

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

func (s *Server) loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	reply, err := handler(ctx, req)

	s.log.Info(
		"[gRPC]",
		s.log.String("method", info.FullMethod),
		s.log.Duration("duration", time.Since(start)),
	)

	return reply, err
}
