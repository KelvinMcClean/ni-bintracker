package gateway

import (
	"bintracker/internal/bintracker"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type APIClient struct {
	BaseURL    *url.URL
	HTTPClient *http.Client
}

func NewAPIClient(baseURL string) (*APIClient, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	return &APIClient{
		BaseURL: parsedURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}, nil
}

func GetBins(cfg bintracker.Config) Response {
	client, err := NewAPIClient("https://ardsandnorthdownbincalendar.azurewebsites.net/api")
	if err != nil {
		panic(err)
	}
	payload, err := client.getBins(cfg)
	if err != nil {
		panic(err)
	}
	return payload
}

func (c *APIClient) getBins(cfg bintracker.Config) (Response, error) {
	// Construct the request URL
	endpoint := c.BaseURL.JoinPath("collectiondates", fmt.Sprintf("%d", cfg.House.Id))
	query := endpoint.Query()
	endpoint.RawQuery = query.Encode()

	req, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)

	if err != nil {
		return Response{}, fmt.Errorf("failed to create request: %w", err)
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return Response{}, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	fmt.Printf("API responded with status: %d\n", resp.StatusCode)
	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return Response{}, fmt.Errorf("failed to decode response: %w", err)
	}
	return response, nil
}
