package mocks

import "delivery-service/storage/cache"

type CacheMock struct {
	GetCampaignsByCountryMock func(string) ([]cache.Campaign, error)
}

func (c CacheMock) GetCampaignsByCountry(country string) ([]cache.Campaign, error) {
	return c.GetCampaignsByCountryMock(country)
}
