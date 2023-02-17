package gateway_api

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/internal/gateway/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

const (
	getUsers = "GetUsers"
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
	expectedResult                  []map[string]interface{}
	expectedStatusCode              int
	expectedError                   map[string]interface{}
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
			expectedResult: []map[string]interface{}{
				{"first_name": "John", "last_name": "Doe"},
			},
		},
		{
			name: "error - wrong credentials",
			arguments: &getUsersHandlerTestCaseArguments{
				authUser:     "test",
				authPassword: "test",
			},
			expectedStatusCode: 401,
			expectedError:      map[string]interface{}{"error": "Unauthorized"},
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

			if testCase.expectedStatusCode == http.StatusOK {
				var jsonResponse []map[string]interface{}
				json.NewDecoder(response.Body).Decode(&jsonResponse)
				suite.Equal(testCase.expectedResult, jsonResponse)
			} else if testCase.expectedStatusCode != http.StatusOK {
				var jsonResponse map[string]interface{}
				json.NewDecoder(response.Body).Decode(&jsonResponse)
				suite.Equal(testCase.expectedError, jsonResponse)
			}

		})
	}
}

func TestGatewayAPITestSuite(t *testing.T) {
	suite.Run(t, new(gatewayAPITestSuite))
}
