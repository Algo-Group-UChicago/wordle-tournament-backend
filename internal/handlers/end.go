package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	. "wordle-tournament-backend/internal/common"
	"wordle-tournament-backend/internal/storage"
	"wordle-tournament-backend/internal/wordle/corpus"
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

// CalculateScore calculates the weighted score from a list of GameState entries.
// Returns the average of (weight * num_guesses) for all solved games.
func calculateScore(games []storage.GameState) float64 {
	if len(games) == 0 {
		return 0.0
	}
	sumOfWeights := 0.0
	totalScore := 0.0
	for _, game := range games {
		weight := corpus.GetWordWeight(game.Answer)
		totalScore += weight * float64(game.NumGuesses)
		sumOfWeights += weight
	}

	return totalScore / sumOfWeights
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
		LogWarning("end", err.Error(), http.StatusBadRequest)
		http.Error(w, "Invalid json body", http.StatusBadRequest)
		return
	}

	if req.TeamID == "" {
		LogWarning("end", errors.New("team_id cannot be empty").Error(), http.StatusBadRequest)
		http.Error(w, "team_id cannot be empty", http.StatusBadRequest)
		return
	}

	if req.RunID == "" {
		LogWarning("end", errors.New("run_id cannot be empty").Error(), http.StatusBadRequest, slog.String("team_id", req.TeamID))
		http.Error(w, "run_id cannot be empty", http.StatusBadRequest)
		return
	}

	LogInfo("end", "entry", slog.String("team_id", req.TeamID), slog.String("run_id", req.RunID))

	// Query ActiveRuns database
	activeRun, err := storage.GetActiveRun(req.TeamID, req.RunID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "expired or not found") {
			statusCode = http.StatusBadRequest
		}

		// StatusCode < 500 differentiates between client and server errors
		if statusCode < 500 {
			LogWarning("end", err.Error(), statusCode, slog.String("team_id", req.TeamID), slog.String("run_id", req.RunID))
		} else {
			LogError("end", err.Error(), statusCode, slog.String("team_id", req.TeamID), slog.String("run_id", req.RunID))
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

	var score, avg float64
	var scorePtr, avgPtr *float64
	var solved bool

	if allSolved {
		score = calculateScore(activeRun.Games)
		avg = totalGuesses / float64(NumTargetWords)
		scorePtr = &score
		avgPtr = &avg
		solved = true
	} else {
		score = 0.0
		avg = 0.0
		scorePtr = nil
		avgPtr = nil
		solved = false
	}

	completedRun := storage.CompletedRun{
		RunID:          req.RunID,
		Score:          scorePtr,
		AverageGuesses: avgPtr,
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
			LogError("end", err.Error(), http.StatusInternalServerError, slog.String("team_id", req.TeamID), slog.String("run_id", req.RunID))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		scoreItem.CompletedRuns = append(scoreItem.CompletedRuns, completedRun)
	}

	if err := storage.PutScore(scoreItem); err != nil {
		LogError("end", err.Error(), http.StatusInternalServerError, slog.String("team_id", req.TeamID), slog.String("run_id", req.RunID))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := storage.RemoveActiveRun(req.TeamID, req.RunID); err != nil {
		LogError("end", err.Error(), http.StatusInternalServerError, slog.String("team_id", req.TeamID), slog.String("run_id", req.RunID))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := EndResponse{
		Score:          score,
		AverageGuesses: avg,
		Solved:         solved,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
