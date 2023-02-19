package users

import (
	"github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/internal/gateway/models"
	usersApi "github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/internal/gateway/proto/users-api"
)

func mapGRPCUsersToUsers(grpCUsers []*usersApi.User) []*models.User {
	users := make([]*models.User, 0, len(grpCUsers))
	for _, grpCUser := range grpCUsers {
		users = append(users, mapGRPCUserToUser(grpCUser))
	}

	return users
}

func mapGRPCUserToUser(grpcUser *usersApi.User) *models.User {
	return &models.User{
		FirstName:   grpcUser.FirstName,
		LastName:    grpcUser.LastName,
		Birthday:    grpcUser.Birthday,
		Address:     grpcUser.Address,
		PhoneNumber: grpcUser.PhoneNumber,
	}
}
