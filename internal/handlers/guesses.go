package handlers

import (
	"encoding/json"
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
			LogWarning("GuessesInvalidRequest", &LogData{Msg: "method not allowed"})
			http.Error(w, "HTTP Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// Potential Issues:
// - If the team_id + run_id are invalid, request returns 500 error when we should return something more helpful.
// - No server-side validation on NumGuesses being less than MAX_GUESSSES (already in middleware)

// /api/guesses is hit once a team wants to submit a guess.
// This handler will:
// 1. Validate the request
// 2. Query the ActiveRuns database to get the run
// 3. Grade the guesses
// 4. Update the ActiveRuns database
// 5. Return the response
func handlePostGuesses(w http.ResponseWriter, r *http.Request) {
	// TODO: uppercase guesses will FAIL
	var req GuessesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		LogWarning("GuessesInvalidRequest", &LogData{Msg: err.Error()})
		http.Error(w, "Invalid json body", http.StatusBadRequest)
		return
	}

	if req.TeamID == "" {
		LogWarning("GuessesEmptyTeamID", &LogData{})
		http.Error(w, "team_id cannot be empty", http.StatusBadRequest)
		return
	}

	if req.RunID == "" {
		LogWarning("GuessesEmptyRunID", &LogData{TeamID: req.TeamID})
		http.Error(w, "run_id cannot be empty", http.StatusBadRequest)
		return
	}

	LogInfo("GuessesProcessRequest", &LogData{TeamID: req.TeamID, RunID: req.RunID})

	if err := wordle.ValidateGuesses(req.Guesses); err != nil {
		LogWarning("GuessesInvalidGuesses", &LogData{TeamID: req.TeamID, RunID: req.RunID, Msg: err.Error()})
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Query the ActiveRuns database to get the current run state
	activeRun, err := storage.GetActiveRun(req.TeamID, req.RunID)
	if err != nil {
		isClientError := strings.Contains(err.Error(), "expired or not found")
		if isClientError {
			LogWarning("GuessesInvalidIDs", &LogData{TeamID: req.TeamID, RunID: req.RunID, Msg: err.Error()})
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			LogError("GuessesInternalError", &LogData{TeamID: req.TeamID, RunID: req.RunID, Msg: err.Error()})
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Extract answers from activeRun.Games
	answers := make([]string, len(activeRun.Games))
	for i := range activeRun.Games {
		answers[i] = activeRun.Games[i].Answer
	}

	// Grade the guesses then update associated metadata (solved, num_guesses, etc)
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

	// Update the ActiveRuns database with the new run state
	if err := storage.PutActiveRun(activeRun); err != nil {
		LogError("GuessesUpdateActiveRunFailure", &LogData{TeamID: req.TeamID, RunID: req.RunID, Msg: err.Error()})
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
