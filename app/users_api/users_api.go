package users_api

import (
	"context"
	proto "github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/app/users_api/proto/users-api"
	"github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/internal/users"
	"github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/internal/users/models"
	file2 "github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/internal/users/repositories/file"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
)

type usersService interface {
	GetUsers(context context.Context) ([]*models.User, error)
}

type GrpcServer struct {
	proto.UnimplementedUsersAPIServiceServer
	Server       *grpc.Server
	usersService usersService
	config       *Config
}

func (g *GrpcServer) GetUsers(ctx context.Context, req *proto.GetUsersRequest) (*proto.GetUsersResponse, error) {
	users, err := g.usersService.GetUsers(ctx)
	if err != nil {
		log.Errorf("failed to get users, error: %v", err)

		return nil, status.Error(codes.Internal, "failed to get users")
	}

	return &proto.GetUsersResponse{
		Users: mapUsersToGRPC(users),
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
	proto.RegisterUsersAPIServiceServer(grpcServer, g)
	g.Server = grpcServer

	g.newServices()
}

func (g *GrpcServer) newServices() {
	repository := file2.New(&file2.Config{
		FileDirectory: g.config.RepositoryFileDirectory,
	})

	g.usersService = users.New(repository)
}

func New(config *Config) *GrpcServer {
	return &GrpcServer{
		config: config,
	}
}
