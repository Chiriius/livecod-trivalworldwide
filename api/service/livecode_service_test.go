package service

import (
	"context"
	"errors"
	"livecode_tribalworldwide/api/entities"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGetUsersService(t *testing.T) {
	testScenarios := []struct {
		testName       string
		mockRepo       *repositoryMock
		mockResponse   []entities.User
		mockError      error
		configureMock  func(*repositoryMock, []entities.User, error)
		expectedOutput []entities.User
		expectedError  error
	}{
		{
			testName: "Successful user retrieval",
			mockRepo: &repositoryMock{},
			mockResponse: []entities.User{
				{UUID: "1", FirstName: "Miguel", LastName: "Gonzalez", Email: "Miguel@gmail.com"},
			},
			mockError: nil,
			configureMock: func(m *repositoryMock, mockResponse []entities.User, mockError error) {
				m.On("GetUsers").Return(mockResponse, mockError)
			},
			expectedOutput: []entities.User{
				{UUID: "1", FirstName: "Miguel", LastName: "Gonzalez", Email: "Miguel@gmail.com"},
			},
			expectedError: nil,
		},
		{
			testName:     "Repository returns error",
			mockRepo:     &repositoryMock{},
			mockResponse: nil,
			mockError:    errors.New("repository error"),
			configureMock: func(m *repositoryMock, mockResponse []entities.User, mockError error) {
				m.On("GetUsers").Return(mockResponse, mockError)
			},
			expectedOutput: nil,
			expectedError:  errors.New("error fetching users"),
		},
	}

	for _, tt := range testScenarios {
		t.Run(tt.testName, func(t *testing.T) {
			if tt.configureMock != nil {
				tt.configureMock(tt.mockRepo, tt.mockResponse, tt.mockError)
			}

			service := NewUserService(tt.mockRepo, logrus.StandardLogger(), context.Background())
			result, err := service.GetUsers()

			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedOutput, result)
		})
	}
}
