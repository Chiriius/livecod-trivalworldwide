package main

import (
	"context"
	"livecode_tribalworldwide/api/server"
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()
	logger := logrus.StandardLogger()
	logger.SetFormatter(&logrus.JSONFormatter{})

	httpAddr := os.Getenv("SERVER_PORT_HTTP")
	if httpAddr == "" {
		httpAddr = ":8080"
	}

	srv, err := server.New(logger, httpAddr, ctx)
	if err != nil {
		logger.Panic("Layer: main ", "Failed to create server:", err)
	}

	srv.Start()
}
