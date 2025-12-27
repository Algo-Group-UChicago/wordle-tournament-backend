package handlers

import (
	"encoding/json"
	"net/http"
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

func handlePostGuesses(w http.ResponseWriter, r *http.Request) {
	// TODO: uppercase guesses will FAIL
	var req GuessesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid json body", http.StatusBadRequest)
		return
	}

	if err := wordle.ValidateTeamId(req.TeamId); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if err := wordle.ValidateGuesses(req.Guesses); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// No validation on whether team_id + run_id are valid
	// TODO: Add team_id + run_id validation. Illegal values currently return a 500 error.
	activeRun, err := storage.GetActiveRun(req.TeamId, req.RunId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
