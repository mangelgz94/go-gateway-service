package gateway

import (
	"context"
	"testing"

	"github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/internal/gateway/models"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

const (
	getUsers           = "GetUsers"
	findNumberPosition = "FindNumberPosition"
)

var (
	emptyContext  = context.Background()
	internalError = errors.New("internal error")
)

type gatewayTestSuite struct {
	suite.Suite
}

type fakeUsersProvider struct {
	mock.Mock
}

func (f *fakeUsersProvider) GetUsers(ctx context.Context) ([]*models.User, error) {
	args := f.Called(ctx)

	return args.Get(0).([]*models.User), args.Error(1)
}

func (f *fakeUsersProvider) Shutdown() {

}

type fakeFindNumberPositionProvider struct {
	mock.Mock
}

func (f *fakeFindNumberPositionProvider) FindNumberPosition(ctx context.Context, number int) (int, error) {
	args := f.Called(ctx, number)

	return args.Get(0).(int), args.Error(1)
}

func (f *fakeFindNumberPositionProvider) Shutdown() {

}

type getUsersTestCase struct {
	name                           string
	argument                       context.Context
	usersProviderGetUsersMockSetup *usersProviderGetUsersMockSetup
	expectedResult                 []*models.User
	expectedError                  error
}

type usersProviderGetUsersMockSetup struct {
	argument       context.Context
	expectedResult []*models.User
	expectedError  error
}

func (suite *gatewayTestSuite) TestGetUsers() {
	testCases := []*getUsersTestCase{
		{
			name:     "success - get users successfully",
			argument: emptyContext,
			usersProviderGetUsersMockSetup: &usersProviderGetUsersMockSetup{
				argument: emptyContext,
				expectedResult: []*models.User{
					{
						FirstName: "John",
						LastName:  "Doe",
					},
				},
			},
			expectedResult: []*models.User{
				{
					FirstName: "John",
					LastName:  "Doe",
				},
			},
		},
		{
			name:     "error - get users",
			argument: emptyContext,
			usersProviderGetUsersMockSetup: &usersProviderGetUsersMockSetup{
				argument:      emptyContext,
				expectedError: internalError,
			},
			expectedError: errors.Wrap(internalError, "usersProvider GetUsers"),
		},
	}

	for _, testCase := range testCases {
		suiteT := suite.T()
		suite.Run(testCase.name, func() {
			usersProvider := &fakeUsersProvider{}
			mockSetup := testCase.usersProviderGetUsersMockSetup
			if mockSetup != nil {
				usersProvider.On(getUsers, mockSetup.argument).Return(mockSetup.expectedResult, mockSetup.expectedError)
			}
			gatewayService := New(usersProvider, nil)
			users, err := gatewayService.GetUsers(testCase.argument)

			usersProvider.AssertExpectations(suiteT)
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

type findNumberPositionTestCase struct {
	name                                                  string
	arguments                                             *findNumberPositionTestCaseArguments
	findNumberPositionProviderFindNumberPositionMockSetup *findNumberPositionProviderFindNumberPositionMockSetup
	expectedResult                                        int
	expectedError                                         error
}

type findNumberPositionProviderFindNumberPositionMockSetup struct {
	arguments      *findNumberPositionMethodArguments
	expectedResult int
	expectedError  error
}

type findNumberPositionMethodArguments struct {
	ctx    context.Context
	number int
}
type findNumberPositionTestCaseArguments struct {
	ctx    context.Context
	number int
}

func (suite *gatewayTestSuite) TestFindNumberPosition() {
	testCases := []*findNumberPositionTestCase{
		{
			name: "success - find number position successfully",
			arguments: &findNumberPositionTestCaseArguments{
				ctx:    emptyContext,
				number: 1,
			},
			findNumberPositionProviderFindNumberPositionMockSetup: &findNumberPositionProviderFindNumberPositionMockSetup{
				arguments: &findNumberPositionMethodArguments{
					ctx:    emptyContext,
					number: 1,
				},
				expectedResult: 1,
			},
			expectedResult: 1,
		},
		{
			name: "error - find number position error",
			arguments: &findNumberPositionTestCaseArguments{
				ctx:    emptyContext,
				number: 1,
			},
			findNumberPositionProviderFindNumberPositionMockSetup: &findNumberPositionProviderFindNumberPositionMockSetup{
				arguments: &findNumberPositionMethodArguments{
					ctx:    emptyContext,
					number: 1,
				},
				expectedError: internalError,
			},
			expectedError: errors.Wrap(internalError, "findNumberPositionProvider FindNumberPosition"),
		},
	}

	for _, testCase := range testCases {
		suiteT := suite.T()
		suite.Run(testCase.name, func() {
			findNumberPositionProvider := &fakeFindNumberPositionProvider{}
			mockSetup := testCase.findNumberPositionProviderFindNumberPositionMockSetup
			if mockSetup != nil {
				findNumberPositionProvider.On(findNumberPosition, mockSetup.arguments.ctx, mockSetup.arguments.number).Return(mockSetup.expectedResult, mockSetup.expectedError)
			}

			gatewayService := New(nil, findNumberPositionProvider)
			numberPosition, err := gatewayService.FindNumberPosition(testCase.arguments.ctx, testCase.arguments.number)

			findNumberPositionProvider.AssertExpectations(suiteT)
			suite.Equal(testCase.expectedResult, numberPosition)
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

func TestGatewayTestSuite(t *testing.T) {
	suite.Run(t, new(gatewayTestSuite))
}
