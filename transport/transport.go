package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"delivery-service/endpoints"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

const (
	getCampaignsUrl = "/v1/delivery"
)

var (
	logger              log.Logger
	ErrMissingParams    = errors.New("missing required parameters: app, country, or os")
	ErrMethodNotAllowed = errors.New("method not allowed")
)

func init() {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "ts", log.DefaultTimestamp, "package", "transport")

	// Set log level debug
	logger = level.NewFilter(logger, level.AllowDebug())

}

// DecodeGetCampaignsRequest decodes the incoming HTTP request into our request struct
func DecodeGetCampaignsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	switch r.Method {
	case "GET":
		level.Info(logger).Log("api", "REQUEST", "method", "GetCampaignsRequest", "url", r.URL.String(), "httpMethod", r.Method)
		break
	default:
		level.Info(logger).Log("api", "REQUEST", "method", "GetCampaignsRequest", "url", r.URL.String(), "httpMethod", r.Method, "err", "Method Not Allowed")
		return nil, ErrMethodNotAllowed
	}

	var request endpoints.GetCampaignsRequest

	if app := r.URL.Query().Get("app"); app != "" {
		request.App = app
	} else {
		level.Error(logger).Log("api", "REQUEST", "method", "DecodeGetCampaignsRequest", "err", "Missing required app parameter")
		return nil, ErrMissingParams
	}

	if country := r.URL.Query().Get("country"); country != "" {
		request.Country = country
	} else {
		level.Error(logger).Log("api", "REQUEST", "method", "DecodeGetCampaignsRequest", "err", "Missing required country parameter")
		return nil, ErrMissingParams
	}

	if os := r.URL.Query().Get("os"); os != "" {
		request.Os = r.URL.Query().Get("os")
	} else {
		level.Error(logger).Log("api", "REQUEST", "method", "DecodeGetCampaignsRequest", "err", "Missing required os parameter")
		return nil, ErrMissingParams
	}

	return request, nil
}

// EncodeResponse encodes the outgoing response as JSON
func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.WriteHeader(http.StatusOK)
	level.Info(logger).Log("api", "RESPONSE", "method", "GetCampaignsRequest", "httpStatusCode", http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

// EncodeErrorResponse encodes the error response and sets the appropriate HTTP status code
func EncodeErrorResponse(_ context.Context, err error, w http.ResponseWriter) {
	if errors.Is(err, ErrMissingParams) {
		w.WriteHeader(http.StatusBadRequest)
	} else if errors.Is(err, ErrMethodNotAllowed) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}

// NewHTTPHandler creates an HTTP handler
func NewHTTPHandler(ep endpoint.Endpoint) http.Handler {
	getCampaignsHandler := httptransport.NewServer(
		ep,
		DecodeGetCampaignsRequest,
		EncodeResponse,
		httptransport.ServerErrorEncoder(EncodeErrorResponse),
	)

	mux := http.NewServeMux()
	mux.Handle(getCampaignsUrl, getCampaignsHandler)
	return mux
}
