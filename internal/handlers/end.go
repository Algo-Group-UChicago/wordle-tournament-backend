package handlers

import (
	"encoding/json"
	"math"
	"net/http"
	"strings"
	"time"

	"wordle-tournament-backend/internal/common"
	"wordle-tournament-backend/internal/storage"
	"wordle-tournament-backend/internal/wordle"
)

type EndRequest struct {
	TeamID string `json:"team_id"`
	RunID  string `json:"run_id"`
}

type EndResponse struct {
	TotalScore     float64 `json:"total_score"`
	AverageGuesses float64 `json:"average_guesses"`
	SolvedCount    int     `json:"solved_count"`
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

	if err := wordle.ValidateTeamId(req.TeamID); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if req.RunID == "" {
		http.Error(w, "run_id cannot be empty", http.StatusBadRequest)
		return
	}

	activeRun, err := storage.GetActiveRun(req.TeamID, req.RunID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "expired or not found") {
			statusCode = http.StatusBadRequest
		}
		http.Error(w, err.Error(), statusCode)
		return
	}

	// Calculate score from game states
	totalGuesses := 0.0
	solvedCount := 0
	for _, game := range activeRun.Games {
		if game.Solved {
			totalGuesses += float64(game.NumGuesses)
			solvedCount++
		}
	}

	// Check if all games are solved
	if solvedCount != common.NumTargetWords {
		http.Error(w, "not all games are solved", http.StatusBadRequest)
		return
	}

	averageGuesses := totalGuesses / float64(common.NumTargetWords)
	totalScore := averageGuesses

	// Create CompletedRun
	completedRun := storage.CompletedRun{
		RunID:          req.RunID,
		TotalScore:     totalScore,
		AverageGuesses: averageGuesses,
		SolvedCount:    solvedCount,
		CompletedAt:    time.Now(),
	}

	// Get or create ScoreItem for the team
	scoreItem, err := storage.GetScore(req.TeamID)
	if err != nil {
		// Team doesn't exist in Scores table, create new entry
		if strings.Contains(err.Error(), "score not found") {
			scoreItem = &storage.ScoreItem{
				TeamID:        req.TeamID,
				BestScore:     totalScore,
				CompletedRuns: []storage.CompletedRun{completedRun},
			}
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// Team exists, append the completed run and update BestScore
		scoreItem.CompletedRuns = append(scoreItem.CompletedRuns, completedRun)
		if totalScore < scoreItem.BestScore || math.Abs(scoreItem.BestScore) < 0.0001 {
			scoreItem.BestScore = totalScore
		}
	}

	// Save to Scores table
	if err := storage.PutScore(scoreItem); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Remove the active run
	if err := storage.RemoveActiveRun(req.TeamID, req.RunID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := EndResponse{
		TotalScore:     totalScore,
		AverageGuesses: averageGuesses,
		SolvedCount:    solvedCount,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
