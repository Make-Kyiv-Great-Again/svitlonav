package handler

import (
	"encoding/json"
	"net/http"

	"svitlonav/internal/service"
	"svitlonav/internal/valhalla"
)

type RouteHandler struct {
	svc *service.RouteService
}

func NewRouteHandler(svc *service.RouteService) *RouteHandler {
	return &RouteHandler{svc: svc}
}

type routeRequestDTO struct {
	From valhalla.LatLon `json:"from"`
	To   valhalla.LatLon `json:"to"`
}

func (h *RouteHandler) GetRoute(w http.ResponseWriter, r *http.Request) {
	var req routeRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	result, err := h.svc.GetLitRoute(r.Context(), req.From, req.To)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *RouteHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("API is running!"))
}
