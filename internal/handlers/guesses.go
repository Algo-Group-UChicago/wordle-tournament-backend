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
	// might be prudent to move validation + grading calls
	// to another function that only takes a GuessesRequest

	// TODO: uppercase guesses will FAIL.
	// what is the expected behavior
	var req GuessesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid json body", http.StatusBadRequest)
		return
	}

	// if err := wordle.ValidateTeamId(req.TeamId); err != nil {
	// 	http.Error(w, err.Error(), http.StatusUnauthorized)
	// 	return
	// }

	// if err := wordle.ValidateGuesses(req.Guesses); err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	response := GuessesResponse{
		Hints: wordle.GradeGuesses(req.Guesses, []string{"afoul", "after"}),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
