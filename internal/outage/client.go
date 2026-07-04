package outage

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)
// outages won't work
type Client struct {
	baseURL string
	http    *http.Client
}

func NewDTEKClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		http:    &http.Client{},
	}
}

func (c *Client) ActiveOutageGroups(ctx context.Context) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/api/v1/outages/active", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("dtek api request: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Groups []string `json:"active_groups"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("dtek api decode: %w", err)
	}

	return result.Groups, nil
}
