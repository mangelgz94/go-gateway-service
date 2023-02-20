//go:build end_to_end_test
// +build end_to_end_test

package end_to_end_tests

import (
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/suite"
)

type endToEndTestSuite struct {
	suite.Suite
}

var (
	baseServiceURL = loadBaseServiceURL()
	authUser       = loadAuthUser()
	authPassword   = loadAuthPassword()
	arraySize      = loadArraySize()
)

func loadArraySize() int {
	arraySizeParameter := os.Getenv("ARRAY_SIZE")
	if arraySizeParameter == "" {
		return 100
	}
	arraySize, err := strconv.Atoi(arraySizeParameter)
	if err != nil {
		return 100
	}

	return arraySize
}

func loadAuthPassword() string {
	authPassword := os.Getenv("AUTH_PASSWORD")
	if authPassword == "" {
		return "password"
	}

	return authPassword
}

func loadAuthUser() string {
	authUser := os.Getenv("AUTH_USER")
	if authUser == "" {
		return "admin"
	}

	return authUser
}

func loadBaseServiceURL() string {
	serviceURL := os.Getenv("BASE_SERVICE_URL")
	if serviceURL == "" {
		return "http://localhost:8090"
	}
	return serviceURL
}

type getUsersTestCase struct {
	name               string
	arguments          *getUsersTestCaseArguments
	expectedStatusCode int
	expectedResult     *getUsersResponse
	expectedError      error
}

type getUsersTestCaseArguments struct {
	authUser     string
	authPassword string
}

func (suite *endToEndTestSuite) TestGetUsers() {
	testCases := []*getUsersTestCase{
		{
			name: "success - get users successfully",
			arguments: &getUsersTestCaseArguments{
				authUser:     authUser,
				authPassword: authPassword,
			},
			expectedResult: &getUsersResponse{
				Users: []*User{
					{
						FirstName:   "John",
						LastName:    "Doe",
						Birthday:    "01-01-2000",
						Address:     "His Address",
						PhoneNumber: "1234567890",
					},
					{
						FirstName:   "Jane",
						LastName:    "Doe",
						Birthday:    "01-01-2001",
						Address:     "Her Address",
						PhoneNumber: "0987654321",
					},
				},
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "error - wrong credentials",
			arguments: &getUsersTestCaseArguments{
				authUser:     "test",
				authPassword: "test",
			},
			expectedResult: &getUsersResponse{
				Error: "unauthorized",
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
	}
	for _, testCase := range testCases {
		suiteT := suite.T()
		suite.Run(testCase.name, func() {
			apiTester := httpexpect.Default(suiteT, baseServiceURL)
			var getUsersResponse getUsersResponse
			response := apiTester.GET("/users").
				WithBasicAuth(testCase.arguments.authUser, testCase.arguments.authPassword).
				Expect().
				Status(testCase.expectedStatusCode).
				JSON().
				Object()
			response.Decode(&getUsersResponse).IsEqual(testCase.expectedResult)
		})
	}
}

type findNumberPositionTestCase struct {
	name               string
	arguments          *findNumberPositionTestCaseArguments
	expectedStatusCode int
	expectedResult     *findNumberPositionResponse
	expectedError      error
}

type findNumberPositionTestCaseArguments struct {
	authUser     string
	authPassword string
	number       *string
}

func (suite *endToEndTestSuite) TestFindNumberPosition() {
	number := "50"
	numberNotFound := strconv.Itoa(loadArraySize() + 1)
	invalidNumber := "invalid"
	emptyNumber := ""
	testCases := []*findNumberPositionTestCase{
		{
			name: "success - find number position",
			arguments: &findNumberPositionTestCaseArguments{
				authUser:     authUser,
				authPassword: authPassword,
				number:       &number,
			},
			expectedStatusCode: http.StatusOK,
			expectedResult:     &findNumberPositionResponse{Number: 50},
		},
		{
			name: "success - number not found",
			arguments: &findNumberPositionTestCaseArguments{
				authUser:     authUser,
				authPassword: authPassword,
				number:       &numberNotFound,
			},
			expectedStatusCode: http.StatusOK,
			expectedResult:     &findNumberPositionResponse{Number: -1},
		},
		{
			name: "error - wrong credentials",
			arguments: &findNumberPositionTestCaseArguments{
				authUser:     "test",
				authPassword: "test",
				number:       &number,
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedResult:     &findNumberPositionResponse{Error: "unauthorized"},
		},
		{
			name: "success - invalid number",
			arguments: &findNumberPositionTestCaseArguments{
				authUser:     authUser,
				authPassword: authPassword,
				number:       &invalidNumber,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResult:     &findNumberPositionResponse{Error: "invalid number"},
		},
		{
			name: "success - empty number",
			arguments: &findNumberPositionTestCaseArguments{
				authUser:     authUser,
				authPassword: authPassword,
				number:       &emptyNumber,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResult:     &findNumberPositionResponse{Error: "invalid number"},
		},
	}
	for _, testCase := range testCases {
		suiteT := suite.T()
		suite.Run(testCase.name, func() {
			apiTester := httpexpect.Default(suiteT, baseServiceURL)
			var findNumberPositionResponse findNumberPositionResponse
			response := apiTester.GET("/find_number_position")
			if testCase.arguments != nil && testCase.arguments.number != nil {
				if testCase.arguments.number != nil {
					response = response.WithQuery("number", *testCase.arguments.number)
				}
			}

			response.WithBasicAuth(testCase.arguments.authUser, testCase.arguments.authPassword).
				Expect().
				Status(testCase.expectedStatusCode).
				JSON().
				Object().Decode(&findNumberPositionResponse).IsEqual(testCase.expectedResult)
		})
	}
}

func TestEndToEndTestSuite(t *testing.T) {
	suite.Run(t, new(endToEndTestSuite))
}
