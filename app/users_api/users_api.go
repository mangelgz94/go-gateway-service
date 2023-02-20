package users_api

import (
	"context"
	"os"
	"path/filepath"
	"time"

	proto "github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/app/users_api/proto/users-api"
	"github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/internal/users"
	"github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/internal/users/models"
	"github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/internal/users/repositories/file"
	"github.com/pkg/errors"
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

func (g *GrpcServer) Start() error {
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

	err := g.newServices()
	if err != nil {
		return errors.Wrap(err, "newServices")
	}

	return nil
}

func (g *GrpcServer) newServices() error {
	wd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "os Getwd")
	}

	fileDirectory, err := filepath.Abs(g.config.RepositoryFileDirectory)
	if err != nil {
		return errors.Wrap(err, "filepath Abs")
	}

	fileDirectory = wd + fileDirectory
	log.Infof("file directory %s", fileDirectory)
	repository := file.New(&file.Config{
		FileDirectory: fileDirectory,
	})

	g.usersService = users.New(repository)

	return nil
}

func New(config *Config) *GrpcServer {
	return &GrpcServer{
		config: config,
	}
}
