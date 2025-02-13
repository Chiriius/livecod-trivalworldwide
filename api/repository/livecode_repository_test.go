package repository

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var mockAPIResponse = `{"results":[{"gender":"male","name":{"first":"Miguel","last":"Gonzalez"},"email":"miguel@gmail.com","location":{"city":"riosucio","country":"Colombia"},"login":{"uuid":"12345"}}]}`

func TestFetchUsers_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, mockAPIResponse)
	}))
	defer server.Close()

	ApiUrl = server.URL

	logger := logrus.New()
	repo := NewLivecodeRepository(logger)
	users, err := repo.GetUsers()

	assert.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, "12345", users[0].UUID)
}
