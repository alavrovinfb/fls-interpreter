package server

import (
	"time"

	"github.com/alavrovinfb/fls-interpreter/pkg/pb"
	"github.com/alavrovinfb/fls-interpreter/pkg/svc"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

func NewGRPCServer(logger *logrus.Logger) (*grpc.Server, error) {
	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(
			keepalive.ServerParameters{
				Time:    time.Duration(viper.GetInt("config.keepalive.time")) * time.Second,
				Timeout: time.Duration(viper.GetInt("config.keepalive.timeout")) * time.Second,
			},
		),
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				// logging middleware
				grpc_logrus.UnaryServerInterceptor(logrus.NewEntry(logger)),

				// Metrics middleware
				grpc_prometheus.UnaryServerInterceptor,
			),
		),
	)

	// register service implementation with the grpcServer
	s, err := svc.NewBasicServer()
	if err != nil {
		return nil, err
	}
	pb.RegisterFlsInterpreterServer(grpcServer, s)
	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)

	return grpcServer, nil
}
