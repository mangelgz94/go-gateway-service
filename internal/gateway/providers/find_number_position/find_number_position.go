package find_number_position

import (
	"context"

	findNumberPositionApi "github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/internal/gateway/proto/find-number-position-api"
	"github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/internal/gateway/providers"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type findNumberPositionAPIConnector interface {
	FindNumberPosition(ctx context.Context, in *findNumberPositionApi.FindNumberPositionRequest, opts ...grpc.CallOption) (*findNumberPositionApi.FindNumberPositionResponse, error)
}

type FindNumberPositionProvider struct {
	grpClient                      *providers.GrpcClient
	findNumberPositionAPIConnector findNumberPositionAPIConnector
}

func (f *FindNumberPositionProvider) FindNumberPosition(ctx context.Context, number int) (int, error) {
	numberResponse, err := f.findNumberPositionAPIConnector.FindNumberPosition(ctx, &findNumberPositionApi.FindNumberPositionRequest{Number: int64(number)})
	if err != nil {
		return 0, errors.Wrap(err, "findNumberPositionAPIConnector FindNumberPosition")
	}

	return int(numberResponse.Position), err
}

func New(findNumberPositionAPIConnector findNumberPositionAPIConnector, grpcClient *providers.GrpcClient) *FindNumberPositionProvider {
	return &FindNumberPositionProvider{
		findNumberPositionAPIConnector: findNumberPositionAPIConnector,
		grpClient:                      grpcClient,
	}
}
