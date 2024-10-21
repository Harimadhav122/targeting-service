package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"delivery-service/endpoints"
	"delivery-service/service"
	"delivery-service/transport"

	"github.com/stretchr/testify/assert"
)

// test 200 http status code
func TestMain1(t *testing.T) {
	// Set up the service, endpoint, and HTTP handler
	svc := service.NewService()
	ep := endpoints.MakeGetCampaignsEndpoint(svc)
	handler := transport.NewHTTPHandler(ep)

	// Create a test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// Make a valid request to the test server
	resp, err := http.Get(server.URL + "/v1/delivery?app=com.gametion.ludokinggame&country=us&os=android")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// test 400 http status code by missing os param
func TestMain2(t *testing.T) {
	// Set up the service, endpoint, and HTTP handler
	svc := service.NewService()
	ep := endpoints.MakeGetCampaignsEndpoint(svc)
	handler := transport.NewHTTPHandler(ep)

	// Create a test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// Make a valid request to the test server
	resp, err := http.Get(server.URL + "/v1/delivery?app=com.gametion.ludokinggame&country=us")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// test 400 http status code by missing country param
func TestMain3(t *testing.T) {
	// Set up the service, endpoint, and HTTP handler
	svc := service.NewService()
	ep := endpoints.MakeGetCampaignsEndpoint(svc)
	handler := transport.NewHTTPHandler(ep)

	// Create a test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// Make a valid request to the test server
	resp, err := http.Get(server.URL + "/v1/delivery?app=com.gametion.ludokinggame&os=android")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// test 400 http status code by missing app param
func TestMain4(t *testing.T) {
	// Set up the service, endpoint, and HTTP handler
	svc := service.NewService()
	ep := endpoints.MakeGetCampaignsEndpoint(svc)
	handler := transport.NewHTTPHandler(ep)

	// Create a test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// Make a valid request to the test server
	resp, err := http.Get(server.URL + "/v1/delivery?country=us&os=ios")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// test 405 http status code by unsupported method
func TestMain5(t *testing.T) {
	// Set up the service, endpoint, and HTTP handler
	svc := service.NewService()
	ep := endpoints.MakeGetCampaignsEndpoint(svc)
	handler := transport.NewHTTPHandler(ep)

	// Create a test server
	server := httptest.NewServer(handler)
	defer server.Close()

	data := "{}"
	reader := bytes.NewReader([]byte(data))

	// Make a valid request to the test server
	resp, err := http.Post(server.URL+"/v1/delivery?app=com.gametion.ludokinggame&country=us", "application/json", reader)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}
