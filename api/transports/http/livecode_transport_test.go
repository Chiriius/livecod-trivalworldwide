package transport

import (
	"encoding/json"
	"livecode_tribalworldwide/api/endpoints"
	"livecode_tribalworldwide/api/entities"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewHTTPHandler(t *testing.T) {
	logger := logrus.New()
	mockEndpoint := new(mockEndpoint)
	mockUsers := []entities.User{
		{UUID: "1", FirstName: "Miguel", LastName: "Gonzalez", Email: "Miguel@gmail.com"},
		{UUID: "2", FirstName: "Miguel", LastName: "Gonzalez", Email: "Miguel@gmail.com"},
	}

	mockEndpoint.On("GetUsers", mock.Anything, mock.Anything).Return(mockUsers, nil)

	endpoints := endpoints.Endpoints{
		GetUsers: mockEndpoint.GetUsers,
	}
	handler := NewHTTPHandler(endpoints, logger)

	t.Run("Test GetUsers Endpoint", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/users", nil)
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		var response []entities.User
		assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
		assert.Len(t, response, len(mockUsers))
	})
}
