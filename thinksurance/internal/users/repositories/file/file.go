package file

import (
	"context"
	"encoding/json"
	"github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/thinksurance/internal/users/models"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"io/fs"
	"os"
	"path/filepath"
)

type FileRepository struct {
	config *Config
}

func (f *FileRepository) GetUsers(context context.Context) ([]*models.User, error) {
	var files []string
	err := filepath.Walk(f.config.FileDirectory, func(path string, info fs.FileInfo, err error) error {
		if filepath.Ext(path) == ".json" {
			files = append(files, path)
		}

		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "filepath Walk")
	}

	var errGroup errgroup.Group
	usersChannel := make(chan *models.User, len(files))
	for _, file := range files {
		file := file
		errGroup.Go(func() error {
			err := f.getUser(file, usersChannel)
			if err != nil {
				return errors.Wrap(err, "getUser")
			}

			return nil
		})
	}

	if err := errGroup.Wait(); err != nil {
		return nil, err
	}

	users := make([]*models.User, 0, len(files))
	for i := 0; i < len(files); i++ {
		users = append(users, <-usersChannel)
	}

	return users, nil

}

func (f *FileRepository) getUser(file string, channel chan *models.User) error {
	userFile, err := os.ReadFile(file)
	if err != nil {
		return errors.Wrap(err, "os ReadFile")
	}

	var user *models.User
	err = json.Unmarshal(userFile, &user)
	if err != nil {
		return errors.Wrap(err, "json Unmarshal")
	}

	channel <- user

	return nil
}

func New(config *Config) *FileRepository {
	return &FileRepository{
		config: config,
	}
}
