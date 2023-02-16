package users_api

import (
	"context"
	proto "github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/app/users_api/proto/users-api"
	"github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/internal/users/models"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	getUsers = "GetUsers"
)

var (
	emptyContext  = context.Background()
	internalError = errors.New("internal error")
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
	argument       context.Context
	expectedResult []*models.User
	expectedError  error
}

func (suite *usersAPITestSuite) TestGetUsers() {
	testCases := []*getUsersTestCase{
		{
			name: "success - get users",
			arguments: &getUsersTestCaseArguments{
				ctx:     emptyContext,
				request: &proto.GetUsersRequest{},
			},
			usersServiceGetUsersMockSetup: &usersServiceGetUsersMockSetup{
				argument: emptyContext,
				expectedResult: []*models.User{
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
			expectedResult: &proto.GetUsersResponse{
				Users: []*proto.User{
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
		},
		{
			name: "error - get users",
			arguments: &getUsersTestCaseArguments{
				ctx:     emptyContext,
				request: &proto.GetUsersRequest{},
			},
			usersServiceGetUsersMockSetup: &usersServiceGetUsersMockSetup{
				argument:      emptyContext,
				expectedError: internalError,
			},
			expectedError: status.Error(codes.Internal, "failed to get users"),
		},
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
				if testCase.expectedError != nil {
					suite.EqualError(err, testCase.expectedError.Error())
				} else if testCase.expectedError == nil {
					suite.Errorf(err, "unexpected error returned")
				}
			}
		})
	}
}

func TestUsersAPITestSuite(t *testing.T) {
	suite.Run(t, new(usersAPITestSuite))
}
