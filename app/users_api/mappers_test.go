package users_api

import (
	proto "github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/app/users_api/proto/users-api"
	"github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/internal/users/models"
	"testing"

	"github.com/stretchr/testify/suite"
)

type mappersTestSuite struct {
	suite.Suite
}

type mapUsersToGRPCTestCase struct {
	name           string
	argument       []*models.User
	expectedResult []*proto.User
}

func (suite *mappersTestSuite) TestMapUsersToGRPC() {
	testCases := []*mapUsersToGRPCTestCase{
		{
			name: "success - map users to grpc",
			argument: []*models.User{
				{
					FirstName:   "John",
					LastName:    "Doe",
					Address:     "My Address",
					Birthday:    "01-01-2000",
					PhoneNumber: "+123456789",
				},
				{
					FirstName:   "Jane",
					LastName:    "Doe",
					Address:     "Not My Address",
					Birthday:    "01-01-2001",
					PhoneNumber: "+987654321",
				},
			},
			expectedResult: []*proto.User{
				{
					FirstName:   "John",
					LastName:    "Doe",
					Address:     "My Address",
					Birthday:    "01-01-2000",
					PhoneNumber: "+123456789",
				},
				{
					FirstName:   "Jane",
					LastName:    "Doe",
					Address:     "Not My Address",
					Birthday:    "01-01-2001",
					PhoneNumber: "+987654321",
				},
			},
		},
		{
			name:           "success - empty users",
			expectedResult: []*proto.User{},
		},
		{
			name: "success - map nil user",
			argument: []*models.User{
				nil,
			},
			expectedResult: []*proto.User{
				nil,
			},
		},
	}

	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			users := mapUsersToGRPC(testCase.argument)
			suite.Equal(testCase.expectedResult, users)
		})
	}
}

func TestMappersTestSuite(t *testing.T) {
	suite.Run(t, new(mappersTestSuite))
}
