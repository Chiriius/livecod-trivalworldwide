package endpoints

import (
	"context"
	"errors"
	"livecode_tribalworldwide/api/entities"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestMakeGetUsersEndpoint(t *testing.T) {
	logger := logrus.New()
	mockSvc := new(mockService)
	ep := MakeGetUsersEndpoint(mockSvc, logger)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mockUsers := []entities.User{{UUID: "1", FirstName: "Migue", LastName: "Gonzalez"}}
		mockSvc.On("GetUsers").Return(mockUsers, nil).Once()

		resp, err := ep(ctx, struct{}{})
		assert.NoError(t, err)
		assert.Equal(t, GetUsersResponse{Users: mockUsers, Err: ""}, resp)

		mockSvc.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockSvc.On("GetUsers").Return(nil, errors.New("fetch error")).Once()

		resp, err := ep(ctx, struct{}{})
		assert.NoError(t, err)
		assert.Equal(t, GetUsersResponse{Users: nil, Err: "fetch error"}, resp)

		mockSvc.AssertExpectations(t)
	})
}
