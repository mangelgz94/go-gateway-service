package file

import (
	"context"

	"github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/thinksurance/internal/services/users/models"
)

type FileRepository struct {
	config *Config
}

func (f *FileRepository) GetUsers(context context.Context) ([]*models.User, error) {
	return nil, nil
}

func New(config *Config) *FileRepository {
	return &FileRepository{
		config: config,
	}
}
