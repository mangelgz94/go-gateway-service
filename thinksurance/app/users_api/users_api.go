package users_api

import (
	"context"

	proto "github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/thinksurance/app/users_api/proto/users-api"
	"github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/thinksurance/internal/services/users/models"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type usersService interface {
	GetUsers(context context.Context) ([]*models.User, error)
}

type GrpcServer struct {
	proto.UnimplementedUsersAPIServiceServer
	Server       grpc.Server
	usersService usersService
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

func New(usersService usersService) *GrpcServer {
	return &GrpcServer{
		usersService: usersService,
	}
}
