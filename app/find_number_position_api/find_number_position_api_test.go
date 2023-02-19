package find_number_position_api

import (
	"context"
	"github.com/stretchr/testify/mock"
	"testing"

	proto "github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/app/find_number_position_api/proto/find-number-position-api"
	"github.com/stretchr/testify/suite"
)

var emptyContext = context.Background()

var (
	findNumberPosition = "FindNumberPosition"
)

type findNumberPositionAPITestSuite struct {
	suite.Suite
}

type fakeFindNumberPositionService struct {
	mock.Mock
}

func (f *fakeFindNumberPositionService) FindNumberPosition(context context.Context, number int) int {
	args := f.Called(context, number)

	return args.Get(0).(int)
}

type findNumberPositionTestCase struct {
	name                                                 string
	arguments                                            *findNumberPositionTestCaseArguments
	findNumberPositionServiceFindNumberPositionMockSetup *findNumberPositionServiceFindNumberPositionMockSetup
	expectedResult                                       *proto.FindNumberPositionResponse
}

type findNumberPositionServiceFindNumberPositionMockSetup struct {
	arguments      *findNumberPositionMethodArguments
	expectedResult int
}

type findNumberPositionMethodArguments struct {
	ctx    context.Context
	number int
}

func (suite *findNumberPositionAPITestSuite) TestFindNumberPosition() {
	testCases := []*findNumberPositionTestCase{
		{
			name: "success - find number position",
			arguments: &findNumberPositionTestCaseArguments{
				ctx:     emptyContext,
				request: &proto.FindNumberPositionRequest{Number: 1},
			},
			findNumberPositionServiceFindNumberPositionMockSetup: &findNumberPositionServiceFindNumberPositionMockSetup{
				arguments: &findNumberPositionMethodArguments{
					ctx:    emptyContext,
					number: 1,
				},
				expectedResult: 1,
			},
			expectedResult: &proto.FindNumberPositionResponse{
				Position: 1,
			},
		},
	}

	for _, testCase := range testCases {
		suiteT := suite.T()
		suite.Run(testCase.name, func() {
			findNumberPositionService := &fakeFindNumberPositionService{}
			mockSetup := testCase.findNumberPositionServiceFindNumberPositionMockSetup
			if mockSetup != nil {
				findNumberPositionService.On(findNumberPosition, mockSetup.arguments.ctx, mockSetup.arguments.number).Return(mockSetup.expectedResult)
			}

			grpcServer := &GrpcServer{findNumberPositionService: findNumberPositionService}
			response, _ := grpcServer.FindNumberPosition(testCase.arguments.ctx, testCase.arguments.request)
			findNumberPositionService.AssertExpectations(suiteT)
			suite.Equal(testCase.expectedResult, response)
		})

	}
}

type findNumberPositionTestCaseArguments struct {
	ctx     context.Context
	request *proto.FindNumberPositionRequest
}

func TestFindNumberPositionAPITestSuite(t *testing.T) {
	suite.Run(t, new(findNumberPositionAPITestSuite))
}
