package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	. "wordle-tournament-backend/internal/common"
	"wordle-tournament-backend/internal/storage"
	"wordle-tournament-backend/internal/wordle"
)

type StartRequest struct {
	TeamID string `json:"team_id"`
}

type StartResponse struct {
	RunID string `json:"run_id"`
}

func StartHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlePostStart(w, r)
		default:
			LogWarning("start", errors.New("method not allowed").Error(), http.StatusMethodNotAllowed)
			http.Error(w, "HTTP Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func handlePostStart(w http.ResponseWriter, r *http.Request) {
	var req StartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		LogWarning("start", err.Error(), http.StatusBadRequest)
		http.Error(w, "Invalid json body", http.StatusBadRequest)
		return
	}

	if err := wordle.ValidateTeamId(req.TeamID); err != nil {
		LogWarning("start", err.Error(), http.StatusUnauthorized, slog.String("team_id", req.TeamID))
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	runID := uuid.New().String()
	LogInfo("start", "entry", slog.String("team_id", req.TeamID), slog.String("run_id", runID))
	if err := storage.PutDefaultActiveRun(req.TeamID, runID); err != nil {
		LogError("start", err.Error(), http.StatusInternalServerError, slog.String("team_id", req.TeamID), slog.String("run_id", runID))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(StartResponse{RunID: runID})
}
