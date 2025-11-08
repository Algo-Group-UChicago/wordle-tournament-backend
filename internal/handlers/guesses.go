package handlers

import (
	"encoding/json"
	"net/http"
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
	// extract GuessesRequest obj from json paylod
	var req GuessesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if !validateTeamId(req.TeamId) {
		http.Error(w, "Invalid team_id", http.StatusUnauthorized)
		return
	}

	if !validateGuesses(req.Guesses) {

	}

	// Process the guesses (placeholder logic for now)
	// TODO: Implement actual Wordle hint generation
	results := make([]GuessResult, len(req.Guesses))
	correctCount := 0

	for i, guess := range req.Guesses {
		// Placeholder: all guesses marked as incorrect with "absent" hints
		results[i] = GuessResult{
			Index:   i,
			Guess:   guess,
			Hints:   []string{"absent", "absent", "absent", "absent", "absent"},
			Correct: false,
		}
	}

	response := GuessResponse{
		GameKey: req.GameKey,
		Results: results,
		Summary: GuessSummary{
			Total:   len(req.Guesses),
			Correct: correctCount,
			Failed:  len(req.Guesses) - correctCount,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func validateTeamId(teamid string) bool {
	return true
}

func validateGuesses(guesses []string) bool {

}

// gradeGuess compares a guess against the answer and returns hints for each letter position.
// Returns a slice of strings where each element is one of: "correct", "present", or "absent"
func gradeGuess(guess, answer string) []string {
	if len(guess) != 5 || len(answer) != 5 {
		return []string{"absent", "absent", "absent", "absent", "absent"}
	}

	hints := make([]string, 5)
	answerRunes := []rune(answer)
	guessRunes := []rune(guess)

	// Track which answer letters have been used
	used := make([]bool, 5)

	// First pass: mark correct positions
	for i := 0; i < 5; i++ {
		if guessRunes[i] == answerRunes[i] {
			hints[i] = "correct"
			used[i] = true
		}
	}

	// Second pass: mark present letters
	for i := 0; i < 5; i++ {
		if hints[i] == "correct" {
			continue
		}

		found := false
		for j := 0; j < 5; j++ {
			if !used[j] && guessRunes[i] == answerRunes[j] {
				hints[i] = "present"
				used[j] = true
				found = true
				break
			}
		}

		if !found {
			hints[i] = "absent"
		}
	}

	return hints
}

