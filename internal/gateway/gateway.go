package gateway

import (
	"context"

	"github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/internal/gateway/models"
	"github.com/pkg/errors"
)

type usersProvider interface {
	GetUsers(ctx context.Context) ([]*models.User, error)
	Shutdown()
}

type findNumberPositionProvider interface {
	FindNumberPosition(ctx context.Context, number int) (int, error)
	Shutdown()
}

type GatewayService struct {
	usersProvider              usersProvider
	findNumberPositionProvider findNumberPositionProvider
}

func (g *GatewayService) Shutdown() {
	if g.usersProvider != nil {
		g.usersProvider.Shutdown()
	}

	if g.findNumberPositionProvider != nil {
		g.findNumberPositionProvider.Shutdown()
	}
}

func (g *GatewayService) FindNumberPosition(ctx context.Context, number int) (int, error) {
	numberPosition, err := g.findNumberPositionProvider.FindNumberPosition(ctx, number)
	if err != nil {
		return 0, errors.Wrap(err, "findNumberPositionProvider FindNumberPosition")
	}

	return numberPosition, nil
}

func (g *GatewayService) GetUsers(ctx context.Context) ([]*models.User, error) {
	users, err := g.usersProvider.GetUsers(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "usersProvider GetUsers")
	}

	return users, nil
}

func New(usersProvider usersProvider, findNumberPositionProvider findNumberPositionProvider) *GatewayService {
	return &GatewayService{
		usersProvider:              usersProvider,
		findNumberPositionProvider: findNumberPositionProvider,
	}
}
