package service

import (
	"context"
	"fmt"

	"svitlonav/internal/repository"
	"svitlonav/internal/valhalla"
)

const preferIlluminatedUseLit = 0.6
const lampDisplayRadiusM = 25

type LampPoint struct {
	ID  int64   `json:"id"`
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type RouteResult struct {
	Route [][2]float64 `json:"route"`
	Lamps []LampPoint  `json:"lamps"`
}

type RouteService struct {
	lampRepo *repository.LampRepository
	valhalla *valhalla.Client
}

func NewRouteService(lampRepo *repository.LampRepository, v *valhalla.Client) *RouteService {
	return &RouteService{
		lampRepo: lampRepo,
		valhalla: v,
	}
}

func (s *RouteService) GetLitRoute(ctx context.Context, from, to valhalla.LatLon) (*RouteResult, error) {
	route, err := s.valhalla.GetPedestrianRoute(ctx, from, to, preferIlluminatedUseLit)
	if err != nil {
		return nil, fmt.Errorf("valhalla route: %w", err)
	}

	nearby, err := s.lampRepo.LampsNearLine(ctx, route, lampDisplayRadiusM)
	if err != nil {
		return nil, fmt.Errorf("lamps near route: %w", err)
	}

	lamps := make([]LampPoint, 0, len(nearby))
	for _, l := range nearby {
		lamps = append(lamps, LampPoint{ID: l.ID, Lat: l.Lat, Lon: l.Lon})
	}

	return &RouteResult{Route: route, Lamps: lamps}, nil
}
