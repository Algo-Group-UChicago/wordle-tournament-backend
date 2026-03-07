package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	. "wordle-tournament-backend/internal/common"
	"wordle-tournament-backend/internal/storage"
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
			LogWarning("StartInvalidRequest", &LogData{Msg: "method not allowed"})
			http.Error(w, "HTTP Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// /api/start is hit once a team wants to start a new run.
// This handler will:
// 1. Validate the request
// 2. Create a new active run for the team
// 3. Return the response
func handlePostStart(w http.ResponseWriter, r *http.Request) {
	var req StartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		LogWarning("StartInvalidRequest", &LogData{Msg: err.Error()})
		http.Error(w, "Invalid json body", http.StatusBadRequest)
		return
	}

	if req.TeamID == "" {
		LogWarning("StartEmptyTeamID", &LogData{})
		http.Error(w, "team_id cannot be empty", http.StatusBadRequest)
		return
	}

	runID := uuid.New().String()
	LogInfo("StartProcessRequest", &LogData{TeamID: req.TeamID, RunID: runID})

	// Create a new active run for the team
	if err := storage.PutDefaultActiveRun(req.TeamID, runID); err != nil {
		LogError("StartCreateActiveRunFailure", &LogData{TeamID: req.TeamID, RunID: runID, Msg: err.Error()})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(StartResponse{RunID: runID})
}
