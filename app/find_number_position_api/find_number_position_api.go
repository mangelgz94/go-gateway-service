package find_number_position_api

import (
	"context"
	"time"

	proto "github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/app/find_number_position_api/proto/find-number-position-api"
	"github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/internal/find_number_position"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type findNumberPositionService interface {
	FindNumberPosition(context context.Context, number int) int
}

type GrpcServer struct {
	proto.UnimplementedFindNumberPositionAPIServiceServer
	Server                    *grpc.Server
	findNumberPositionService findNumberPositionService
	config                    *Config
}

func (g *GrpcServer) GetNumberPosition(ctx context.Context, req *proto.GetNumberPositionRequest) (*proto.GetNumberPositionResponse, error) {
	return &proto.GetNumberPositionResponse{
		Position: int64(g.findNumberPositionService.FindNumberPosition(ctx, int(req.Number))),
	}, nil
}

func (g *GrpcServer) Start() {
	keepAliveEnforcementPolicy := keepalive.EnforcementPolicy{
		MinTime:             time.Duration(g.config.ServerKeepAliveEnforcementMinTime) * time.Second,
		PermitWithoutStream: g.config.ServerKeepAlivePermitWithoutStream,
	}

	keepAliveServerParameters := keepalive.ServerParameters{
		MaxConnectionIdle:     time.Duration(g.config.ServerKeepAliveMaxConnectionIdle) * time.Second,
		MaxConnectionAge:      time.Duration(g.config.ServerKeepAliveMaxConnectionAge) * time.Second,
		MaxConnectionAgeGrace: time.Duration(g.config.ServerKeepAliveMaxConnectionAgeGrace) * time.Second,
		Time:                  time.Duration(g.config.ServerKeepAliveTime) * time.Second,
		Timeout:               time.Duration(g.config.ServerKeepAliveTimeout) * time.Second,
	}

	grpcServer := grpc.NewServer(grpc.KeepaliveEnforcementPolicy(keepAliveEnforcementPolicy), grpc.KeepaliveParams(keepAliveServerParameters))
	proto.RegisterFindNumberPositionAPIServiceServer(grpcServer, g)
	g.Server = grpcServer

	g.findNumberPositionService = find_number_position.New(&find_number_position.Config{ArraySize: g.config.ArraySize})
}

func New(config *Config) *GrpcServer {
	return &GrpcServer{
		config: config,
	}
}
