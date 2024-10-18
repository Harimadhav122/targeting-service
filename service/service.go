package service

import (
	"context"
)

// Campaign represents a campaign entity
type Campaign struct {
	Cid string `json:"cid"`
	Img string `json:"img"`
	Cta string `json:"cta"`
}

// Service defines the behavior of our campaign service
type Service interface {
	GetCampaigns(ctx context.Context, app string, country string, os string) ([]Campaign, error)
}

// campaignService is the implementation of the Service interface
type campaignService struct {
	campaigns []Campaign
}

// NewService creates and returns a new Campaign Service
func NewService() Service {
	return &campaignService{
		campaigns: []Campaign{
			{Cid: "spotify", Img: "https://somelink", Cta: "Download"},
			{Cid: "duolingo", Img: "https://somelink2", Cta: "Install"},
			{Cid: "subwaysurfer", Img: "https://somelink3", Cta: "Play"},
		},
	}
}

// GetCampaigns implements the business logic
func (s *campaignService) GetCampaigns(ctx context.Context, app, country, os string) ([]Campaign, error) {

	var result []Campaign

	for _, campaign := range s.campaigns {
		if campaign.Cid == "duolingo" && country != "US" && os == "android" {
			result = append(result, campaign)
		} else if campaign.Cid == "spotify" && (country == "US" || country == "Canada") {
			result = append(result, campaign)
		} else if campaign.Cid == "subwaysurfer" && os == "android" && app == "com.gametion.ludokinggame" {
			result = append(result, campaign)
		}
	}

	return result, nil
}
