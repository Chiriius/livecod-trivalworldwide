package endpoints

import (
	"context"
	"livecode_tribalworldwide/api/entities"
	"livecode_tribalworldwide/api/service"

	"github.com/go-kit/kit/endpoint"
	"github.com/sirupsen/logrus"
)

type GetUsersResponse struct {
	Users []entities.User `json:"users"`
	Err   string          `json:"err,omitempty"`
}

type Endpoints struct {
	GetUsers endpoint.Endpoint
}

func MakeServerEndpoints(s service.LiveService, logger logrus.FieldLogger) Endpoints {
	return Endpoints{
		GetUsers: MakeGetUsersEndpoint(s, logger),
	}
}

func MakeGetUsersEndpoint(s service.LiveService, logger logrus.FieldLogger) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		users, err := s.GetUsers()
		if err != nil {
			logger.Errorln("Layer:user_endpoint", "Method:MakeLoginEndpoint", err)
			return GetUsersResponse{Users: nil, Err: err.Error()}, nil
		}
		return GetUsersResponse{Users: users, Err: ""}, nil
	}
}
