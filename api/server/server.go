package server

import (
	"context"
	"livecode_tribalworldwide/api/endpoints"
	"livecode_tribalworldwide/api/repository"
	"livecode_tribalworldwide/api/service"
	transport "livecode_tribalworldwide/api/transports/http"
	"net/http"

	"github.com/sirupsen/logrus"
)

type Server struct {
	httpMux  *http.ServeMux
	httpAddr string
	logger   logrus.FieldLogger
}

func New(logger logrus.FieldLogger, httpAddr string, ctx context.Context) (*Server, error) {
	livecodeRepo := repository.NewLivecodeRepository(logger)
	userService := service.NewUserService(livecodeRepo, logger, ctx)

	livecodeEndpoints := endpoints.MakeServerEndpoints(userService, logger)
	httpHandler := transport.NewHTTPHandler(livecodeEndpoints, logger)

	httpMux := http.NewServeMux()
	httpMux.Handle("/", httpHandler)

	return &Server{
		httpMux:  httpMux,
		httpAddr: httpAddr,
		logger:   logger,
	}, nil
}

func (s *Server) Start() error {
	s.logger.Infof("Server starting on port %s", s.httpAddr)
	err := http.ListenAndServe(s.httpAddr, s.httpMux)
	if err != nil {
		s.logger.Errorf("HTTP server failed: %v", err)
		return err
	}
	return nil
}
