package handlers

import (
	"encoding/json"
	"net/http"
	"wordle-tournament-backend/internal/common"
	"wordle-tournament-backend/internal/wordle"
	"wordle-tournament-backend/internal/wordle/corpus"
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

	// TODO: make answers dependent on a team's id
	possibleAnswers := corpus.GetGradingAnswerKey()
	answers := possibleAnswers[:common.NumTargetWords]

	response := GuessesResponse{
		Hints: wordle.GradeGuesses(req.Guesses, answers),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
