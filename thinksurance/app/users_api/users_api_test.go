package users_api

import (
	"context"
	"testing"

	proto "github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/thinksurance/app/users_api/proto/users-api"
	"github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/thinksurance/internal/services/users/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

const (
	getUsers = "GetUsers"
)

type usersAPITestSuite struct {
	suite.Suite
}

type fakeUsersService struct {
	mock.Mock
}

func (f *fakeUsersService) GetUsers(ctx context.Context) ([]*models.User, error) {
	args := f.Called(ctx)

	return args.Get(0).([]*models.User), args.Error(1)
}

type getUsersTestCase struct {
	name                          string
	arguments                     *getUsersTestCaseArguments
	usersServiceGetUsersMockSetup *usersServiceGetUsersMockSetup
	expectedResult                *proto.GetUsersResponse
	expectedError                 error
}

type getUsersTestCaseArguments struct {
	ctx     context.Context
	request *proto.GetUsersRequest
}

type usersServiceGetUsersMockSetup struct {
	argument       *context.Context
	expectedResult []*models.User
	expectedError  error
}

func (suite *usersAPITestSuite) TestGetUsers() {
	testCases := []*getUsersTestCase{
		{},
	}

	for _, testCase := range testCases {
		suiteT := suite.T()
		suite.Run(testCase.name, func() {
			usersService := &fakeUsersService{}
			if testCase.usersServiceGetUsersMockSetup != nil {
				usersService.On(getUsers, testCase.usersServiceGetUsersMockSetup.argument).Return(testCase.usersServiceGetUsersMockSetup.expectedResult, testCase.usersServiceGetUsersMockSetup.expectedError)
			}

			grpcServer := &GrpcServer{usersService: usersService}
			users, err := grpcServer.GetUsers(testCase.arguments.ctx, testCase.arguments.request)
			usersService.AssertExpectations(suiteT)
			suite.Equal(testCase.expectedResult, users)
			if err != nil {
				suite.EqualError(err, testCase.expectedError.Error())
			}
		})
	}
}

func TestUsersAPITestSuite(t *testing.T) {
	suite.Run(t, new(usersAPITestSuite))
}
