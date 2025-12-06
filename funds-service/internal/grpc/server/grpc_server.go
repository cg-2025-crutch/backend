package server

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/cg-2025-crutch/backend/funds-service/internal/config"
	"github.com/cg-2025-crutch/backend/funds-service/internal/funds/handler"
	"github.com/cg-2025-crutch/backend/funds-service/internal/funds/service"
	pb "github.com/cg-2025-crutch/backend/funds-service/internal/grpc/gen"
	"github.com/cg-2025-crutch/backend/funds-service/internal/infrastructure/log"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	srv  *grpc.Server
	host string
}

func NewGrpcServer(conf config.ServerConfig) *GrpcServer {
	return &GrpcServer{
		srv:  grpc.NewServer(),
		host: conf.GRPCHost + ":" + conf.GRPCPort,
	}
}

func (s *GrpcServer) RegisterGRPC(srvc service.FundsServicer) {
	pb.RegisterFundsServiceServer(s.srv, handler.NewGRPCHandler(srvc))
}

func (s *GrpcServer) Run(ctx context.Context) error {

	g, ctx := errgroup.WithContext(ctx)

	l := log.FromContext(ctx)

	g.Go(func() error {
		lis, err := net.Listen("tcp", s.host)
		if err != nil {
			return fmt.Errorf("grpc server listen error: %w", err)
		}

		l.Infof("grpc server started on: %v", s.host)
		if err := s.srv.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			return fmt.Errorf("grpc server serve error: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		<-ctx.Done()

		l.Info("shutting down grpc server gracefully...")

		ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		done := make(chan struct{})
		go func() {
			s.srv.GracefulStop()
			close(done)
		}()

		select {
		case <-done:
			l.Info("grpc server stopped gracefully")
			return nil
		case <-ctxShutdown.Done():
			l.Warn("forced stop grpc server after 10 seconds")
			s.srv.Stop()
			return ctxShutdown.Err()
		}
	})

	return g.Wait()
}
