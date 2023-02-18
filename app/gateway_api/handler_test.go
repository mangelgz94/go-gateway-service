package gateway_api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/internal/gateway/models"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

const (
	getUsers           = "GetUsers"
	findNumberPosition = "FindNumberPosition"
)

var (
	internalError = errors.New("internal error")
)

type gatewayAPITestSuite struct {
	suite.Suite
}

type fakeGatewayService struct {
	mock.Mock
}

func (f *fakeGatewayService) GetUsers(ctx context.Context) ([]*models.User, error) {
	args := f.Called(ctx)

	return args.Get(0).([]*models.User), args.Error(1)
}

func (f *fakeGatewayService) FindNumberPosition(ctx context.Context, number int) (int, error) {
	args := f.Called(ctx, number)

	return args.Get(0).(int), args.Error(1)
}

type getUsersHandlerTestCase struct {
	name                            string
	arguments                       *getUsersHandlerTestCaseArguments
	gatewayServiceGetUsersMockSetup *gatewayServiceGetUsersMockSetup
	expectedResult                  getUsersResponse
	expectedStatusCode              int
}

type getUsersHandlerTestCaseArguments struct {
	authUser     string
	authPassword string
}

type gatewayServiceGetUsersMockSetup struct {
	expectedResult []*models.User
	expectedError  error
}

func (suite *gatewayAPITestSuite) TestGetUsersHandler() {
	testCases := []*getUsersHandlerTestCase{
		{
			name: "success - get users",
			arguments: &getUsersHandlerTestCaseArguments{
				authUser:     "admin",
				authPassword: "password",
			},
			gatewayServiceGetUsersMockSetup: &gatewayServiceGetUsersMockSetup{
				expectedResult: []*models.User{
					{
						FirstName: "John",
						LastName:  "Doe",
					},
				},
			},
			expectedStatusCode: 200,
			expectedResult: getUsersResponse{
				Users: []*models.User{
					{
						FirstName: "John",
						LastName:  "Doe",
					},
				},
			},
		},
		{
			name: "error - wrong credentials",
			arguments: &getUsersHandlerTestCaseArguments{
				authUser:     "test",
				authPassword: "test",
			},
			expectedStatusCode: 401,
			expectedResult:     getUsersResponse{Error: "unauthorized"},
		},
		{
			name: "error - get users",
			arguments: &getUsersHandlerTestCaseArguments{
				authUser:     "admin",
				authPassword: "password",
			},
			gatewayServiceGetUsersMockSetup: &gatewayServiceGetUsersMockSetup{
				expectedError: internalError,
			},
			expectedStatusCode: 500,
			expectedResult:     getUsersResponse{Error: "internal server error"},
		},
	}

	for _, testCase := range testCases {
		suiteT := suite.T()
		suite.Run(testCase.name, func() {
			gatewayService := &fakeGatewayService{}
			mockSetup := testCase.gatewayServiceGetUsersMockSetup
			if mockSetup != nil {
				gatewayService.On(getUsers, mock.Anything).Return(mockSetup.expectedResult, mockSetup.expectedError)
			}
			gatewayServer := &GatewayServer{gatewayService: gatewayService, config: &Config{AuthUser: "admin", AuthPassword: "password"}}

			handler := gatewayServer.GetUsersHandler()
			server := httptest.NewServer(handler.Middleware.HandleAuthorization(handler.Func))
			req, err := http.NewRequest("GET", server.URL, nil)
			req.SetBasicAuth(testCase.arguments.authUser, testCase.arguments.authPassword)
			if err != nil {
				log.Print(err)
				os.Exit(1)
			}
			response, err := server.Client().Do(req)
			if err != nil {
				fmt.Println("Error when sending request to the server")
				return
			}
			defer response.Body.Close()

			gatewayService.AssertExpectations(suiteT)
			suite.Equal(testCase.expectedStatusCode, response.StatusCode)

			var getUsersResponse getUsersResponse
			json.NewDecoder(response.Body).Decode(&getUsersResponse)
			suite.Equal(testCase.expectedResult, getUsersResponse)

		})
	}
}

type findNumberPositionTestCase struct {
	name                                      string
	arguments                                 *findNumberPositionTestCaseArguments
	gatewayServiceFindNumberPositionMockSetup *gatewayServiceFindNumberPositionMockSetup
	expectedStatusCode                        int
	expectedResult                            findNumberPositionResponse
}

type gatewayServiceFindNumberPositionMockSetup struct {
	argument       int
	expectedResult int
	expectedError  error
}

type findNumberPositionTestCaseArguments struct {
	authUser     string
	authPassword string
	number       *string
}

func (suite *gatewayAPITestSuite) TestFindNumberPosition() {
	number := "1"
	emptyNumber := ""
	testCases := []*findNumberPositionTestCase{
		{
			name: "success - find number position successfully",
			arguments: &findNumberPositionTestCaseArguments{
				authUser:     "admin",
				authPassword: "password",
				number:       &number,
			},
			gatewayServiceFindNumberPositionMockSetup: &gatewayServiceFindNumberPositionMockSetup{
				argument:       1,
				expectedResult: 1,
			},
			expectedStatusCode: http.StatusOK,
			expectedResult:     findNumberPositionResponse{Number: 1},
		},
		{
			name: "error - wrong credentials",
			arguments: &findNumberPositionTestCaseArguments{
				authUser:     "test",
				authPassword: "test",
				number:       &number,
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedResult:     findNumberPositionResponse{Error: "unauthorized"},
		},
		{
			name: "error - empty number",
			arguments: &findNumberPositionTestCaseArguments{
				authUser:     "admin",
				authPassword: "password",
				number:       &emptyNumber,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResult:     findNumberPositionResponse{Error: "invalid number"},
		},
		{
			name: "error - no number",
			arguments: &findNumberPositionTestCaseArguments{
				authUser:     "admin",
				authPassword: "password",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResult:     findNumberPositionResponse{Error: "invalid number"},
		},
	}

	for _, testCase := range testCases {
		suiteT := suite.T()
		suite.Run(testCase.name, func() {
			gatewayService := &fakeGatewayService{}
			mockSetup := testCase.gatewayServiceFindNumberPositionMockSetup
			if mockSetup != nil {
				gatewayService.On(findNumberPosition, mock.Anything, mockSetup.argument).Return(mockSetup.expectedResult, mockSetup.expectedError)
			}

			gatewayServer := &GatewayServer{gatewayService: gatewayService, config: &Config{AuthUser: "admin", AuthPassword: "password"}}

			handler := gatewayServer.FindNumberPositionHandler()
			server := httptest.NewServer(handler.Middleware.HandleAuthorization(handler.Func))
			req, err := http.NewRequest("GET", server.URL, nil)
			req.SetBasicAuth(testCase.arguments.authUser, testCase.arguments.authPassword)
			if err != nil {
				log.Print(err)
				os.Exit(1)
			}

			if testCase.arguments.number != nil {
				q := req.URL.Query()
				q.Add("number", *testCase.arguments.number)
				req.URL.RawQuery = q.Encode()
			}

			response, err := server.Client().Do(req)
			if err != nil {
				fmt.Println("Error when sending request to the server")
				return
			}
			defer response.Body.Close()

			gatewayService.AssertExpectations(suiteT)
			suite.Equal(testCase.expectedStatusCode, response.StatusCode)

			var findNumberPositionResponse findNumberPositionResponse
			json.NewDecoder(response.Body).Decode(&findNumberPositionResponse)
			suite.Equal(testCase.expectedResult, findNumberPositionResponse)
		})
	}
}

func TestGatewayAPITestSuite(t *testing.T) {
	suite.Run(t, new(gatewayAPITestSuite))
}
