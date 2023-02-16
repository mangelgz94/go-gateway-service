package find_number_position

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type findNumberPositionServiceTestSuite struct {
	suite.Suite
}

var (
	emptyContext = context.Background()
)

type getPositionTestCase struct {
	name           string
	arguments      *getPositionTestCaseArguments
	expectedResult int
}

type getPositionTestCaseArguments struct {
	ctx       context.Context
	number    int
	arraySize int
}

func (suite *findNumberPositionServiceTestSuite) TestGetIndex() {
	testCases := []*getPositionTestCase{
		{
			name: "success - find position successfully",
			arguments: &getPositionTestCaseArguments{
				ctx:       emptyContext,
				number:    25,
				arraySize: 100,
			},
			expectedResult: 25,
		},
		{
			name: "failed - position was not found",
			arguments: &getPositionTestCaseArguments{
				ctx:       emptyContext,
				number:    101,
				arraySize: 100,
			},
			expectedResult: -1,
		},
	}

	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			findNumberPositionService := New(&Config{ArraySize: testCase.arguments.arraySize})
			position := findNumberPositionService.FindNumberPosition(testCase.arguments.ctx, testCase.arguments.number)

			suite.Equal(testCase.expectedResult, position)
		})
	}
}

func TestFindNumberPositionServiceTestSuite(t *testing.T) {
	suite.Run(t, new(findNumberPositionServiceTestSuite))
}
