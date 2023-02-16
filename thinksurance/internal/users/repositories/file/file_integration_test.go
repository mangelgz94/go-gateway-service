package file_test

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/thinksurance/internal/users/models"
	"github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/thinksurance/internal/users/repositories/file"
	"github.com/stretchr/testify/suite"
)

type FileRepositoryTestSuite struct {
	suite.Suite
}

var (
	fileDirectory   = loadFileDirectory()
	emptyBackground = context.Background()
)

func loadFileDirectory() string {
	fileDirectory, _ := os.LookupEnv("REPOSITORY_FILE_DIRECTORY")
	if fileDirectory == "" {
		filePath, _ := filepath.Abs("../../../../../users_json")

		return filePath
	}

	return fileDirectory
}

type getUsersTestCase struct {
	name           string
	arguments      *getUsersTestCaseArguments
	expectedResult []*models.User
	expectedError  error
}

type getUsersTestCaseArguments struct {
	ctx   context.Context
	users []string
}

func (suite *FileRepositoryTestSuite) TestGetUsers() {

	fileUsers := make([]string, 0, 100)
	users := make([]*models.User, 0, 100)
	for i := 0; i < 99; i++ {
		fileUsers = append(fileUsers, `
				{
				  "first_name": "Jane`+strconv.Itoa(i)+`",
				  "last_name": "Doe",
				  "birthday": "01-01-2001",
				  "address": "Her Address",
				  "phone_number": "0987654321"
				}`)
		users = append(users, &models.User{
			FirstName:   "Jane" + strconv.Itoa(i),
			LastName:    "Doe",
			Address:     "Her Address",
			Birthday:    "01-01-2001",
			PhoneNumber: "0987654321",
		})
	}
	testCases := []*getUsersTestCase{
		{
			name: "success - get users",
			arguments: &getUsersTestCaseArguments{
				ctx:   emptyBackground,
				users: fileUsers,
			},
			expectedResult: users,
		},
		{
			name: "error - corrupted file",
			arguments: &getUsersTestCaseArguments{
				ctx: emptyBackground,
				users: []string{`corrupted`,
					`{
                  "first_name": "John",
				  "last_name": "Doe",
				  "birthday": "01-01-2000",
				  "address": "His Address",
				  "phone_number": "1234567890"
				}`},
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
	}

	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			fileRepository := file.New(&file.Config{FileDirectory: fileDirectory})
			for index, user := range testCase.arguments.users {

				_ = os.WriteFile(fileDirectory+"/"+strconv.Itoa(index)+".json", []byte(user), 0644)
			}

			users, err := fileRepository.GetUsers(testCase.arguments.ctx)
			suite.ElementsMatch(testCase.expectedResult, users)
			if err != nil {
				if testCase.expectedError != nil {
					suite.EqualError(err, testCase.expectedError.Error())
				} else if testCase.expectedError == nil {
					suite.Errorf(err, "unexpected error returned")
				}
			}

			for index, _ := range testCase.arguments.users {
				_ = os.Remove(fileDirectory + "/" + strconv.Itoa(index) + ".json")
			}

		})
	}
}

func TestFileRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(FileRepositoryTestSuite))
}
