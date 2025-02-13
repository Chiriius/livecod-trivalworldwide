package transport

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type mockEndpoint struct {
	mock.Mock
}

func (m *mockEndpoint) GetUsers(ctx context.Context, request interface{}) (interface{}, error) {
	args := m.Called(ctx, request)
	return args.Get(0), args.Error(1)
}
