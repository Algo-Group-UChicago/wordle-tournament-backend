package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"wordle-tournament-backend/internal/common"
	"wordle-tournament-backend/internal/storage"
	"wordle-tournament-backend/internal/wordle"
)

type GuessesRequest struct {
	TeamId  string   `json:"team_id"`
	RunId   string   `json:"run_id"`
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
		http.Error(w, "Invalid json body", http.StatusBadRequest)
		return
	}

	if req.TeamId == "" {
		http.Error(w, "team_id cannot be empty", http.StatusBadRequest)
		return
	}

	if req.RunId == "" {
		http.Error(w, "run_id cannot be empty", http.StatusBadRequest)
		return
	}

	if err := wordle.ValidateGuesses(req.Guesses); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	activeRun, err := storage.GetActiveRun(req.TeamId, req.RunId)
	if err != nil {
		// Must distinguish between (team_id, run_id) being invalid and network issues causing the request to fail.
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "expired or not found") {
			statusCode = http.StatusBadRequest
		}
		http.Error(w, err.Error(), statusCode)
		return
	}

	// Process guesses and update activeRun GameStates
	answers := make([]string, len(activeRun.Games))
	for i := range req.Guesses {
		answers[i] = activeRun.Games[i].Answer

		if req.Guesses[i] == common.DummyGuess {
			activeRun.Games[i].Solved = true
		} else if !activeRun.Games[i].Solved {
			activeRun.Games[i].NumGuesses++
		}
	}

	if err := storage.PutActiveRun(activeRun); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := GuessesResponse{
		Hints: wordle.GradeGuesses(req.Guesses, answers),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
