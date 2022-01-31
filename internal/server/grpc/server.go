package internalgrpc

import (
	"context"
	"net"

	"github.com/alexMolokov/otus-rotate-banner/api/pb"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcValidator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"google.golang.org/grpc"
)

//go:generate mockgen -destination=grpc/mock_test.go -package=internalgrpc -source=server.go
type Application interface {
	AddBannerToSlot(ctx context.Context, bannerID, slotID int64) error
	RemoveBannerFromSlot(ctx context.Context, bannerID, slotID int64) error
	CountTransition(ctx context.Context, bannerID, slotID, sgID int64) error
	ChooseBanner(ctx context.Context, slotID, sgID int64) (bannerID int64, err error)
}

type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warning(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

type Server struct {
	GRPCAddr   string
	grpcServer *grpc.Server
	service    *RotatorService
}

func NewServer(logger Logger, app Application, grpcAddr string) *Server {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(
			LoggerInterceptor(logger),
			grpcValidator.UnaryServerInterceptor(),
		)),
	)

	return &Server{
		GRPCAddr:   grpcAddr,
		grpcServer: grpcServer,
		service:    &RotatorService{app: app, logger: logger},
	}
}

func (s *Server) Start() error {
	lsn, err := net.Listen("tcp", s.GRPCAddr)
	if err != nil {
		return err
	}
	pb.RegisterRotatorServer(s.grpcServer, s.service)
	return s.grpcServer.Serve(lsn)
}

func (s *Server) Stop(_ context.Context) error {
	s.grpcServer.GracefulStop()
	return nil
}
