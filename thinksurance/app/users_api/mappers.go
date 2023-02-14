package users_api

import (
	proto "github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/thinksurance/app/users_api/proto/users-api"
	"github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/thinksurance/internal/services/users/models"
)

func mapUsersToGRPC(users []*models.User) []*proto.User {
	grpcUsers := make([]*proto.User, 0, len(users))
	for _, user := range users {
		grpcUsers = append(grpcUsers, mapUserToGRPC(user))
	}

	return grpcUsers
}

func mapUserToGRPC(user *models.User) *proto.User {
	if user == nil {
		return nil
	}

	return &proto.User{
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Birthday:    user.Birthday,
		Address:     user.Address,
		PhoneNumber: user.PhoneNumber,
	}
}
