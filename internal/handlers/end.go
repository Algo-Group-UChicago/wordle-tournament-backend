package handlers

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"strings"
	"time"

	"wordle-tournament-backend/internal/common"
	"wordle-tournament-backend/internal/storage"
)

type EndRequest struct {
	TeamID string `json:"team_id"`
	RunID  string `json:"run_id"`
}

type EndResponse struct {
	Score          float64 `json:"score"`
	AverageGuesses float64 `json:"average_guesses"`
	Solved         bool    `json:"solved"`
}

func EndHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlePostEnd(w, r)
		default:
			http.Error(w, "HTTP Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func handlePostEnd(w http.ResponseWriter, r *http.Request) {
	var req EndRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid json body", http.StatusBadRequest)
		return
	}

	if req.TeamID == "" {
		http.Error(w, "team_id cannot be empty", http.StatusBadRequest)
		return
	}

	if req.RunID == "" {
		http.Error(w, "run_id cannot be empty", http.StatusBadRequest)
		return
	}

	// Query ActiveRuns database
	activeRun, err := storage.GetActiveRun(req.TeamID, req.RunID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "expired or not found") {
			statusCode = http.StatusBadRequest
		}
		http.Error(w, err.Error(), statusCode)
		return
	}

	// Check if all entries are solved and calculate values
	allSolved := true
	totalGuesses := 0.0
	for _, game := range activeRun.Games {
		if !game.Solved {
			allSolved = false
			break
		}
		totalGuesses += float64(game.NumGuesses)
	}

	var score, averageGuesses float64
	var solved bool

	if allSolved {
		averageGuesses = totalGuesses / float64(common.NumTargetWords)
		score = storage.CalculateScore(activeRun.Games)
		solved = true
	} else {
		// If not all games are solved, use positive infinity
		score = math.Inf(1)
		averageGuesses = math.Inf(1)
		solved = false
	}

	completedRun := storage.CompletedRun{
		RunID:          req.RunID,
		Score:          score,
		AverageGuesses: averageGuesses,
		Solved:         solved,
		CompletedAt:    time.Now(),
	}

	scoreItem, err := storage.GetScore(req.TeamID)
	if err != nil {
		if errors.Is(err, storage.ErrScoreNotFound) {
			scoreItem = &storage.ScoreItem{
				TeamID:        req.TeamID,
				CompletedRuns: []storage.CompletedRun{completedRun},
			}
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		scoreItem.CompletedRuns = append(scoreItem.CompletedRuns, completedRun)
	}

	if err := storage.PutScore(scoreItem); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := storage.RemoveActiveRun(req.TeamID, req.RunID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := EndResponse{
		Score:          score,
		AverageGuesses: averageGuesses,
		Solved:         solved,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
