package endpoints

import (
	"context"
	"delivery-service/service"
	"os"
	"time"

	"delivery-service/metrics"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

var logger log.Logger

func init() {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "ts", log.DefaultTimestamp, "package", "endpoints")

	// Set log level debug
	logger = level.NewFilter(logger, level.AllowDebug())

}

// GetCampaignsRequest is the struct for incoming request parameters
type GetCampaignsRequest struct {
	Params map[string]string
	Limit  int
	Page   int
}

// GetCampaignsResponse represents the response for the GetCampaigns API
type GetCampaignsResponse struct {
	Campaigns []service.Campaign
}

// MakeGetCampaignsEndpoint creates an endpoint for the GetCampaigns service
func MakeGetCampaignsEndpoint(svc service.Service) endpoint.Endpoint {

	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetCampaignsRequest)
		start := time.Now()

		campaigns, err := svc.GetCampaigns(ctx, req.Params, req.Limit, req.Page)
		if err != nil {
			level.Error(logger).Log("method", "GetCampaignsEndpoint", "err", err, "took", time.Since(start))
			return nil, err
		}

		level.Info(logger).Log("method", "GetCampaignsEndpoint", "took", time.Since(start))
		metrics.HttpRequestLatency.Observe(time.Since(start).Seconds())
		return GetCampaignsResponse{Campaigns: campaigns}, nil
	}
}
