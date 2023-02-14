package users

import (
	"context"

	"github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/thinksurance/internal/users/models"
	"github.com/pkg/errors"
)

type repository interface {
	GetUsers(ctx context.Context) ([]*models.User, error)
}

type UsersService struct {
	repository repository
}

func (u *UsersService) GetUsers(context context.Context) ([]*models.User, error) {
	users, err := u.repository.GetUsers(context)
	if err != nil {

		return nil, errors.Wrap(err, "repository GetUsers")
	}

	return users, nil
}

func New(repository repository) *UsersService {
	return &UsersService{
		repository: repository,
	}
}
