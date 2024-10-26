package transport

import (
	"context"
	local_error "delivery-service/errors"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"delivery-service/endpoints"
	"delivery-service/metrics"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	getCampaignsUrl = "/v1/delivery"
	metricsUrl      = "/metrics"
)

var logger log.Logger

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
		return nil, &local_error.ErrMethodNotAllowed{Method: r.Method}
	}

	request := endpoints.GetCampaignsRequest{Params: make(map[string]string)}

	params := r.URL.Query()

	for key, value := range params {
		if key == "limit" {
			if limit, err := strconv.Atoi(value[0]); err == nil {
				request.Limit = limit
			}
		} else if key == "page" {
			if page, err := strconv.Atoi(value[0]); err == nil {
				request.Page = page
			}
		} else {
			request.Params[key] = value[0]
		}
	}

	if _, ok := request.Params["app"]; !ok {
		level.Error(logger).Log("api", "REQUEST", "method", "DecodeGetCampaignsRequest", "err", "Missing required app parameter")
		return nil, &local_error.ErrMissingParams{Param: "app", Method: r.Method}
	}

	if _, ok := request.Params["country"]; !ok {
		level.Error(logger).Log("api", "REQUEST", "method", "DecodeGetCampaignsRequest", "err", "Missing required country parameter")
		return nil, &local_error.ErrMissingParams{Param: "country", Method: r.Method}
	}

	if _, ok := request.Params["os"]; !ok {
		level.Error(logger).Log("api", "REQUEST", "method", "DecodeGetCampaignsRequest", "err", "Missing required os parameter")
		return nil, &local_error.ErrMissingParams{Param: "os", Method: r.Method}
	}

	if limit := r.URL.Query().Get("limit"); limit == "" {
		level.Error(logger).Log("api", "REQUEST", "method", "DecodeGetCampaignsRequest", "err", "Missing required limit parameter")
		return nil, &local_error.ErrMissingParams{Param: "limit", Method: r.Method}
	}

	if page := r.URL.Query().Get("page"); page == "" {
		level.Error(logger).Log("api", "REQUEST", "method", "DecodeGetCampaignsRequest", "err", "Missing required page parameter")
		return nil, &local_error.ErrMissingParams{Param: "page", Method: r.Method}
	}

	return request, nil
}

// EncodeResponse encodes the outgoing response as JSON
func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	statusCode := http.StatusOK
	w.WriteHeader(statusCode)
	level.Info(logger).Log("api", "RESPONSE", "method", "GetCampaignsRequest", "httpStatusCode", statusCode)
	metrics.HttpRequestCount.With("method", "GET", "code", strconv.Itoa(statusCode)).Add(1)
	return json.NewEncoder(w).Encode(response)
}

// EncodeErrorResponse encodes the error response and sets the appropriate HTTP status code
func EncodeErrorResponse(_ context.Context, err error, w http.ResponseWriter) {
	paramError := err.(local_error.Error)
	metrics.HttpRequestCount.With("method", paramError.GetMethod(), "code", strconv.Itoa(paramError.GetCode())).Add(1)
	w.WriteHeader(paramError.GetCode())
	json.NewEncoder(w).Encode(map[string]string{"error": paramError.Error()})
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
	mux.Handle(metricsUrl, promhttp.Handler())
	return mux
}
