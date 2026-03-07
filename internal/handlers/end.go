package handlers

import (
	"encoding/json"
	"errors"
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
			LogWarning("EndInvalidRequest", &LogData{Msg: "method not allowed"})
			http.Error(w, "HTTP Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// /api/end is hit once a run is complete.
// This handler will:
// 1. Validate the request
// 2. Query the ActiveRuns database to get the run
// 3. Confirm all games are solved
// 4. Calculate the score and average guesses
// 5. Update the Scores database
// 6. Remove the run from the ActiveRuns database
// 7. Return the response
func handlePostEnd(w http.ResponseWriter, r *http.Request) {
	var req EndRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		LogWarning("EndInvalidRequest", &LogData{Msg: err.Error()})
		http.Error(w, "Invalid json body", http.StatusBadRequest)
		return
	}

	if req.TeamID == "" {
		LogWarning("EndEmptyTeamID", &LogData{})
		http.Error(w, "team_id cannot be empty", http.StatusBadRequest)
		return
	}

	if req.RunID == "" {
		LogWarning("EndEmptyRunID", &LogData{TeamID: req.TeamID})
		http.Error(w, "run_id cannot be empty", http.StatusBadRequest)
		return
	}

	LogInfo("EndProcessRequest", &LogData{TeamID: req.TeamID, RunID: req.RunID})

	// Query ActiveRuns database for current run state (should've ended by now)
	activeRun, err := storage.GetActiveRun(req.TeamID, req.RunID)
	if err != nil {
		isClientError := strings.Contains(err.Error(), "expired or not found")
		if isClientError {
			LogWarning("EndInvalidIDs", &LogData{TeamID: req.TeamID, RunID: req.RunID, Msg: err.Error()})
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			LogError("EndGetActiveRunFailure", &LogData{TeamID: req.TeamID, RunID: req.RunID, Msg: err.Error()})
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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

	// Calculate score and average guesses if all games are solved
	// Otherwise, set score and average guesses to 0
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

	// Query Scores database for team's run history
	// Add the completed run to the team's run history
	scoreItem, err := storage.GetScore(req.TeamID)
	if err != nil {
		if errors.Is(err, storage.ErrScoreNotFound) {
			scoreItem = &storage.ScoreItem{
				TeamID:        req.TeamID,
				CompletedRuns: []storage.CompletedRun{completedRun},
			}
		} else {
			LogError("EndGetScoreFailure", &LogData{TeamID: req.TeamID, RunID: req.RunID, Msg: err.Error()})
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		scoreItem.CompletedRuns = append(scoreItem.CompletedRuns, completedRun)
	}

	// Update the Scores database with the new run history
	if err := storage.PutScore(scoreItem); err != nil {
		LogError("EndUpdateScoreFailure", &LogData{TeamID: req.TeamID, RunID: req.RunID, Msg: err.Error()})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Remove the run from the ActiveRuns database
	if err := storage.RemoveActiveRun(req.TeamID, req.RunID); err != nil {
		LogError("EndRemoveActiveRunFailure", &LogData{TeamID: req.TeamID, RunID: req.RunID, Msg: err.Error()})
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
