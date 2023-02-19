package find_number_position

import (
	"context"
	"github.com/pkg/errors"
	"testing"

	findNumberPositionApi "github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/internal/gateway/proto/find-number-position-api"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

const (
	findNumberPosition = "FindNumberPosition"
)

var (
	emptyContext  = context.Background()
	internalError = errors.New("internal error")
)

type findNumberPositionTestSuite struct {
	suite.Suite
}

type fakeFindNumberPositionConnector struct {
	mock.Mock
}

func (f *fakeFindNumberPositionConnector) FindNumberPosition(ctx context.Context, request *findNumberPositionApi.FindNumberPositionRequest, opts ...grpc.CallOption) (*findNumberPositionApi.FindNumberPositionResponse, error) {
	args := f.Called(ctx, request, opts)

	return args.Get(0).(*findNumberPositionApi.FindNumberPositionResponse), args.Error(1)
}

type findNumberPositionTestCase struct {
	name                                                   string
	arguments                                              *findNumberPositionTestCaseArguments
	findNumberPositionConnectorFindNumberPositionMockSetup *findNumberPositionConnectorFindNumberPositionMockSetup
	expectedResult                                         int
	expectedError                                          error
}

type findNumberPositionTestCaseArguments struct {
	ctx    context.Context
	number int
}

type findNumberPositionConnectorFindNumberPositionMockSetup struct {
	arguments      *findNumberPositionMethodArguments
	expectedResult *findNumberPositionApi.FindNumberPositionResponse
	expectedError  error
}

type findNumberPositionMethodArguments struct {
	ctx     context.Context
	request *findNumberPositionApi.FindNumberPositionRequest
}

func (suite *findNumberPositionTestSuite) TestFindNumberPosition() {
	testCases := []*findNumberPositionTestCase{
		{
			name: "success - find number position",
			arguments: &findNumberPositionTestCaseArguments{
				ctx:    emptyContext,
				number: 1,
			},
			findNumberPositionConnectorFindNumberPositionMockSetup: &findNumberPositionConnectorFindNumberPositionMockSetup{
				arguments: &findNumberPositionMethodArguments{
					ctx:     emptyContext,
					request: &findNumberPositionApi.FindNumberPositionRequest{Number: 1},
				},
				expectedResult: &findNumberPositionApi.FindNumberPositionResponse{Position: 1},
			},
			expectedResult: 1,
		},
		{
			name: "error - find number position",
			arguments: &findNumberPositionTestCaseArguments{
				ctx:    emptyContext,
				number: 1,
			},
			findNumberPositionConnectorFindNumberPositionMockSetup: &findNumberPositionConnectorFindNumberPositionMockSetup{
				arguments: &findNumberPositionMethodArguments{
					ctx:     emptyContext,
					request: &findNumberPositionApi.FindNumberPositionRequest{Number: 1},
				},
				expectedError: internalError,
			},
			expectedError: errors.Wrap(internalError, "findNumberPositionAPIConnector FindNumberPosition"),
		},
	}

	for _, testCase := range testCases {
		suiteT := suite.T()
		suite.Run(testCase.name, func() {
			findNumberPositionConnector := &fakeFindNumberPositionConnector{}
			mockSetup := testCase.findNumberPositionConnectorFindNumberPositionMockSetup
			if mockSetup != nil {
				findNumberPositionConnector.On(findNumberPosition, mockSetup.arguments.ctx, mockSetup.arguments.request, []grpc.CallOption(nil)).Return(mockSetup.expectedResult, mockSetup.expectedError)
			}
			findNumberPositionProvider := New(findNumberPositionConnector, nil)
			numberPosition, err := findNumberPositionProvider.FindNumberPosition(testCase.arguments.ctx, testCase.arguments.number)
			findNumberPositionConnector.AssertExpectations(suiteT)
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

func TestFindNumberPositionTestSuite(t *testing.T) {
	suite.Run(t, new(findNumberPositionTestSuite))
}
