package handlers

import (
	"encoding/json"
	"net/http"
	"wordle-tournament-backend/internal/wordle"
)

type GuessesRequest struct {
	TeamId  string   `json:"team_id"`
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
	var req GuessesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid json body", http.StatusBadRequest)
		return
	}

	if err := wordle.ValidateTeamId(req.TeamId); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if err := wordle.ValidateGuessesLength(req.Guesses); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := wordle.ValidateGuesses(req.Guesses); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: Get answers for each guess and grade them
	// For now, return placeholder hints
	hints := make([]string, len(req.Guesses))
	for i := range req.Guesses {
		hints[i] = "absent,absent,absent,absent,absent"
	}

	response := GuessesResponse{
		Hints: hints,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
