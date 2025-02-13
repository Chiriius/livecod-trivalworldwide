package transport

import (
	"context"
	"encoding/json"
	"livecode_tribalworldwide/api/endpoints"
	"net/http"

	httpTransport "github.com/go-kit/kit/transport/http"
	"github.com/sirupsen/logrus"
)

func NewHTTPHandler(endpoint endpoints.Endpoints, logger logrus.FieldLogger) http.Handler {
	m := http.NewServeMux()

	m.Handle("/users", httpTransport.NewServer(
		endpoint.GetUsers,
		decodeRequest,
		encodeResponse,
	))
	return m

}

func decodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return struct{}{}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}
