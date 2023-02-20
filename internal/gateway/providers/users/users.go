package users

import (
	"context"

	"github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/internal/gateway/models"
	usersApi "github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/internal/gateway/proto/users-api"
	"github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/internal/gateway/providers"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type usersAPIConnector interface {
	GetUsers(ctx context.Context, in *usersApi.GetUsersRequest, opts ...grpc.CallOption) (*usersApi.GetUsersResponse, error)
}

type UsersProvider struct {
	grpClient         *providers.GrpcClient
	usersAPIConnector usersAPIConnector
}

func (f *UsersProvider) Shutdown() {
	f.grpClient.CloseConnection()
}

func (u *UsersProvider) GetUsers(ctx context.Context) ([]*models.User, error) {
	usersResponse, err := u.usersAPIConnector.GetUsers(ctx, &usersApi.GetUsersRequest{})
	if err != nil {
		return nil, errors.Wrap(err, "usersAPIConnector GetUsers")
	}

	return mapGRPCUsersToUsers(usersResponse.Users), nil
}

func New(usersAPIConnector usersAPIConnector, grpcClient *providers.GrpcClient) *UsersProvider {
	return &UsersProvider{
		usersAPIConnector: usersAPIConnector,
		grpClient:         grpcClient,
	}
}
