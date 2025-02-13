package endpoints

import (
	"livecode_tribalworldwide/api/entities"

	"github.com/stretchr/testify/mock"
)

type mockService struct {
	mock.Mock
}

func (m *mockService) GetUsers() ([]entities.User, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.User), args.Error(1)
}
