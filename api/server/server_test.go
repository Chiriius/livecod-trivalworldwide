package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewServer(t *testing.T) {
	logger := logrus.New()
	ctx := context.Background()
	httpAddr := ":8080"
	srv, err := New(logger, httpAddr, ctx)
	require.NoError(t, err)
	require.NotNil(t, srv)
	assert.Equal(t, httpAddr, srv.httpAddr)
}

func TestServerStart(t *testing.T) {
	logger := logrus.New()
	ctx := context.Background()
	httpAddr := ":8080"
	srv, err := New(logger, httpAddr, ctx)
	require.NoError(t, err)
	require.NotNil(t, srv)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	srv.httpMux.ServeHTTP(w, req)
	resp := w.Result()

	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
