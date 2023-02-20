package users

import (
	"context"
	"github.com/pkg/errors"
	"testing"

	"github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/internal/gateway/models"
	usersApi "github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/internal/gateway/proto/users-api"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

const (
	getUsers = "GetUsers"
)

var (
	emptyContext  = context.Background()
	internalError = errors.New("internal error")
)

type UsersTestSuite struct {
	suite.Suite
}

type fakeUsersAPIConnector struct {
	mock.Mock
}

func (f *fakeUsersAPIConnector) GetUsers(ctx context.Context, getUsersResponse *usersApi.GetUsersRequest, opts ...grpc.CallOption) (*usersApi.GetUsersResponse, error) {
	args := f.Called(ctx, getUsersResponse, opts)

	return args.Get(0).(*usersApi.GetUsersResponse), args.Error(1)
}

type getUsersTestCase struct {
	name                               string
	argument                           context.Context
	usersAPIConnectorGetUsersMockSetup *usersAPIConnectorGetUsersMockSetup
	expectedResult                     []*models.User
	expectedError                      error
}

type usersAPIConnectorGetUsersMockSetup struct {
	arguments      *getUsersMethodArguments
	expectedResult *usersApi.GetUsersResponse
	expectedError  error
}

type getUsersMethodArguments struct {
	ctx     context.Context
	request *usersApi.GetUsersRequest
}

func (suite *UsersTestSuite) TestGetUsers() {
	testCases := []*getUsersTestCase{
		{
			name:     "success - get users",
			argument: emptyContext,
			usersAPIConnectorGetUsersMockSetup: &usersAPIConnectorGetUsersMockSetup{
				arguments: &getUsersMethodArguments{
					ctx:     emptyContext,
					request: &usersApi.GetUsersRequest{},
				},
				expectedResult: &usersApi.GetUsersResponse{
					Users: []*usersApi.User{
						{
							FirstName:   "John",
							LastName:    "Doe",
							Birthday:    "01-01-2000",
							Address:     "His Address",
							PhoneNumber: "1234567890",
						},
					},
				},
			},
			expectedResult: []*models.User{
				{
					FirstName:   "John",
					LastName:    "Doe",
					Birthday:    "01-01-2000",
					Address:     "His Address",
					PhoneNumber: "1234567890",
				},
			},
		},
		{
			name:     "error - get users",
			argument: emptyContext,
			usersAPIConnectorGetUsersMockSetup: &usersAPIConnectorGetUsersMockSetup{
				arguments: &getUsersMethodArguments{
					ctx:     emptyContext,
					request: &usersApi.GetUsersRequest{},
				},
				expectedError: internalError,
			},
			expectedError: errors.Wrap(internalError, "usersAPIConnector GetUsers"),
		},
	}

	for _, testCase := range testCases {
		suiteT := suite.T()
		suite.Run(testCase.name, func() {
			usersAPIConnector := &fakeUsersAPIConnector{}
			mockSetup := testCase.usersAPIConnectorGetUsersMockSetup
			if mockSetup != nil {
				usersAPIConnector.On(getUsers, mockSetup.arguments.ctx, mockSetup.arguments.request, []grpc.CallOption(nil)).Return(mockSetup.expectedResult, mockSetup.expectedError)
			}
			usersProvider := New(usersAPIConnector, nil)
			users, err := usersProvider.GetUsers(testCase.argument)
			suite.Equal(testCase.expectedResult, users)
			usersAPIConnector.AssertExpectations(suiteT)
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

func TestUsersTestSuite(t *testing.T) {
	suite.Run(t, new(UsersTestSuite))
}
