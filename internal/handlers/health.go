package handlers

import (
	"encoding/json"
	"net/http"
)

type HealthResponse struct {
	Status      string `json:"status"`
	Environment string `json:"environment"`
}

func HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := HealthResponse{
			Status: "healthy",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
