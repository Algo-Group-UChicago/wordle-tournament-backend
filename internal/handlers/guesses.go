package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"wordle-tournament-backend/internal/common"
	. "wordle-tournament-backend/internal/common"
	"wordle-tournament-backend/internal/storage"
	"wordle-tournament-backend/internal/wordle"
)

type GuessesRequest struct {
	TeamID  string   `json:"team_id"`
	RunID   string   `json:"run_id"`
	Guesses []string `json:"guesses"`
}

type GuessesResponse struct {
	Hints []string `json:"hints"`
}

func GuessesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlePostGuesses(w, r)
		default:
			LogWarning("guesses", errors.New("method not allowed").Error(), http.StatusMethodNotAllowed)
			http.Error(w, "HTTP Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// Potential Issues:
// - If the team_id + run_id are invalid, request returns 500 error when we should return something more helpful.
// - No server-side validation on NumGuesses being less than MAX_GUESSSES (already in middleware)
func handlePostGuesses(w http.ResponseWriter, r *http.Request) {
	// TODO: uppercase guesses will FAIL
	var req GuessesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		LogWarning("guesses", err.Error(), http.StatusBadRequest)
		http.Error(w, "Invalid json body", http.StatusBadRequest)
		return
	}

	if req.TeamID == "" {
		LogWarning("guesses", errors.New("team_id cannot be empty").Error(), http.StatusBadRequest)
		http.Error(w, "team_id cannot be empty", http.StatusBadRequest)
		return
	}

	if req.RunID == "" {
		LogWarning("guesses", errors.New("run_id cannot be empty").Error(), http.StatusBadRequest, slog.String("team_id", req.TeamID))
		http.Error(w, "run_id cannot be empty", http.StatusBadRequest)
		return
	}

	LogInfo("guesses", "entry", slog.String("team_id", req.TeamID), slog.String("run_id", req.RunID))

	if err := wordle.ValidateGuesses(req.Guesses); err != nil {
		LogWarning("guesses", err.Error(), http.StatusBadRequest, slog.String("team_id", req.TeamID), slog.String("run_id", req.RunID))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	activeRun, err := storage.GetActiveRun(req.TeamID, req.RunID)
	if err != nil {
		// Must distinguish between (team_id, run_id) being invalid and network issues causing the request to fail.
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "expired or not found") {
			statusCode = http.StatusBadRequest
		}
		if statusCode < 500 {
			LogWarning("guesses", err.Error(), statusCode, slog.String("team_id", req.TeamID), slog.String("run_id", req.RunID))
		} else {
			LogError("guesses", err.Error(), statusCode, slog.String("team_id", req.TeamID), slog.String("run_id", req.RunID))
		}
		http.Error(w, err.Error(), statusCode)
		return
	}

	// Extract answers from activeRun.Games
	answers := make([]string, len(activeRun.Games))
	for i := range activeRun.Games {
		answers[i] = activeRun.Games[i].Answer
	}

	hints := wordle.GradeGuesses(req.Guesses, answers)

	for i, hint := range hints {
		// If the guess is DummyGuess, the game is already solved
		if req.Guesses[i] == common.DummyGuess {
			activeRun.Games[i].Solved = true
		}

		if req.Guesses[i] != common.DummyGuess && !activeRun.Games[i].Solved {
			activeRun.Games[i].NumGuesses++
		}

		if hint == strings.Repeat("O", common.WordLength) {
			activeRun.Games[i].Solved = true
		}
	}

	if err := storage.PutActiveRun(activeRun); err != nil {
		LogError("guesses", err.Error(), http.StatusInternalServerError, slog.String("team_id", req.TeamID), slog.String("run_id", req.RunID))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := GuessesResponse{
		Hints: hints,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
