package main

import (
	"net/http"
	"os"

	"delivery-service/endpoints"
	"delivery-service/service"
	_ "delivery-service/storage"
	"delivery-service/transport"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

func main() {
	// Set up logger
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "ts", log.DefaultTimestamp, "package", "main")

	// Set log level debug
	logger = level.NewFilter(logger, level.AllowDebug())

	// Initialize the service
	svc := service.NewService()

	// Create the endpoint
	getCampaignsEndpoint := endpoints.MakeGetCampaignsEndpoint(svc)

	// Create the HTTP handler
	httpHandler := transport.NewHTTPHandler(getCampaignsEndpoint)

	// Start the HTTP server
	level.Info(logger).Log("msg", "Starting server on port :8080")
	if err := http.ListenAndServe(":8080", httpHandler); err != nil {
		level.Error(logger).Log("msg", "Failed Starting server on port :8080")
	}
}
