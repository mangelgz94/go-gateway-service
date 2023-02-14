package users

import (
	"context"
	"testing"

	"github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/thinksurance/internal/users/models"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

const (
	getUsers = "GetUsers"
)

var (
	emptyContext  = context.Background()
	internalError = errors.New("internal error")
)

type usersServiceTestSuite struct {
	suite.Suite
}

type fakeRepository struct {
	mock.Mock
}

func (f *fakeRepository) GetUsers(ctx context.Context) ([]*models.User, error) {
	args := f.Called(ctx)

	return args.Get(0).([]*models.User), args.Error(1)
}

type getUsersTestCase struct {
	name                        string
	argument                    context.Context
	repositoryGetUsersMockSetup *repositoryGetUsersMockSetup
	expectedResult              []*models.User
	expectedError               error
}

type repositoryGetUsersMockSetup struct {
	argument       context.Context
	expectedResult []*models.User
	expectedError  error
}

func (suite *usersServiceTestSuite) TestGetUsers() {
	testCases := []*getUsersTestCase{
		{
			name:     "success - get users",
			argument: emptyContext,
			repositoryGetUsersMockSetup: &repositoryGetUsersMockSetup{
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
		{
			name:     "error - get users",
			argument: emptyContext,
			repositoryGetUsersMockSetup: &repositoryGetUsersMockSetup{
				argument:      emptyContext,
				expectedError: internalError,
			},
			expectedError: errors.Wrap(internalError, "repository GetUsers"),
		},
	}

	for _, testCase := range testCases {
		suiteT := suite.T()
		suite.Run(testCase.name, func() {
			repository := &fakeRepository{}
			if testCase.repositoryGetUsersMockSetup != nil {
				repository.On(getUsers, testCase.repositoryGetUsersMockSetup.argument).Return(testCase.repositoryGetUsersMockSetup.expectedResult, testCase.repositoryGetUsersMockSetup.expectedError)
			}

			usersService := New(repository)
			users, err := usersService.GetUsers(testCase.argument)
			repository.AssertExpectations(suiteT)
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

func TestUsersServiceTestSuite(t *testing.T) {
	suite.Run(t, new(usersServiceTestSuite))
}
