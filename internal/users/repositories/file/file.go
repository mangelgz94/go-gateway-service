package file

import (
	"context"
	"encoding/json"
	"github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/internal/users/models"
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
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

	var wg sync.WaitGroup
	usersChannel := make(chan *models.User, len(files))
	for _, file := range files {
		wg.Add(1)
		file := file
		go func() {
			defer wg.Done()
			f.getUser(file, usersChannel)
		}()

	}

	wg.Wait()
	close(usersChannel)

	users := make([]*models.User, 0, len(files))
	for user := range usersChannel {
		users = append(users, user)
	}

	return users, nil

}

func (f *FileRepository) getUser(file string, channel chan *models.User) {
	userFile, err := os.ReadFile(file)
	if err != nil {
		log.Errorf("Cannot read file %s, error: %v", file, err)

		return
	}

	var user *models.User
	err = json.Unmarshal(userFile, &user)
	if err != nil {
		log.Errorf("Cannot unmarshal file %s, error: %v", file, err)

		return
	}

	channel <- user
}

func New(config *Config) *FileRepository {
	return &FileRepository{
		config: config,
	}
}
