package valhalla

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	baseURL string
	http    *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		http:    &http.Client{},
	}
}

type LatLon struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type routeRequest struct {
	Locations      []LatLon        `json:"locations"`
	Costing        string          `json:"costing"`
	CostingOptions *costingOptions `json:"costing_options,omitempty"`
}

type costingOptions struct {
	Pedestrian *pedestrianOptions `json:"pedestrian,omitempty"`
}

type pedestrianOptions struct {
	UseLit float64 `json:"use_lit"`
}

type tripResponse struct {
	Trip struct {
		Legs []struct {
			Shape string `json:"shape"`
		} `json:"legs"`
	} `json:"trip"`
}

func (c *Client) GetPedestrianRoute(ctx context.Context, from, to LatLon, useLit float64) ([][2]float64, error) {
	reqBody := routeRequest{
		Locations: []LatLon{from, to},
		Costing:   "pedestrian",
	}
	if useLit > 0 {
		reqBody.CostingOptions = &costingOptions{
			Pedestrian: &pedestrianOptions{UseLit: useLit},
		}
	}

	buf, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/route", bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("valhalla request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("valhalla returned %d: %s", resp.StatusCode, string(body))
	}

	var trip tripResponse
	if err := json.Unmarshal(body, &trip); err != nil {
		return nil, fmt.Errorf("decode valhalla response: %w", err)
	}
	if len(trip.Trip.Legs) == 0 {
		return nil, fmt.Errorf("valhalla returned no legs for %v -> %v", from, to)
	}

	return DecodePolyline(trip.Trip.Legs[0].Shape, 6), nil
}
