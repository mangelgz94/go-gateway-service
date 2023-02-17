package gateway

import (
	"context"

	"github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/internal/gateway/models"
)

type usersProvider interface {
	GetUsers(ctx context.Context) ([]*models.User, error)
}

type findNumberPositionProvider interface {
	FindNumberPosition(ctx context.Context, number int) (int, error)
}

type GatewayService struct {
	usersProvider              usersProvider
	findNumberPositionProvider findNumberPositionProvider
}

func (g *GatewayService) FindNumberPosition(ctx context.Context, number int) (int, error) {
	return 0, nil
}

func (g *GatewayService) GetUsers(ctx context.Context) ([]*models.User, error) {
	return nil, nil
}

func New(usersProvider usersProvider, findNumberPositionProvider findNumberPositionProvider) *GatewayService {
	return &GatewayService{
		usersProvider:              usersProvider,
		findNumberPositionProvider: findNumberPositionProvider,
	}
}
