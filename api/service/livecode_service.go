package service

import (
	"context"
	"errors"
	"livecode_tribalworldwide/api/entities"
	"livecode_tribalworldwide/api/repository"

	"github.com/sirupsen/logrus"
)

type LiveService interface {
	GetUsers() ([]entities.User, error)
}

type liveService struct {
	ctx        context.Context
	repository repository.LivecodeRepository
	logger     logrus.FieldLogger
}

func NewUserService(repo repository.LivecodeRepository, logger logrus.FieldLogger, ctx context.Context) LiveService {
	return &liveService{
		ctx:        ctx,
		repository: repo,
		logger:     logger,
	}
}

func (s *liveService) GetUsers() ([]entities.User, error) {
	users, err := s.repository.GetUsers()
	if err != nil {
		s.logger.Error("Error fetching users from repository: ", err)
		return nil, errors.New("error fetching users")
	}
	return users, nil
}
