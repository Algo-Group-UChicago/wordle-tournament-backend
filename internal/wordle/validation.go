package wordle

import (
	"errors"
	"fmt"
	"wordle-tournament-backend/internal/dictionary"
)

const NumTargetWords = 1000

var (
	ErrInvalidGuessLength = errors.New("invalid number of guesses")
	ErrInvalidWordLength  = errors.New("word must be exactly 5 characters")
	ErrInvalidTeamId      = errors.New("invalid team_id")
)

func ValidateGuessesLength(guesses []string) error {
	if len(guesses) != NumTargetWords {
		return fmt.Errorf("%w: expected %d, got %d", ErrInvalidGuessLength, NumTargetWords, len(guesses))
	}
	return nil
}

func ValidateGuess(guess string) error {
	if len(guess) != WordLength {
		return fmt.Errorf("%w: %s", ErrInvalidWordLength, guess)
	}

	if !dictionary.IsInCorpus(guess) {
		return fmt.Errorf("word not in corpus: %s", guess)
	}

	return nil
}

func ValidateGuesses(guesses []string) error {
	for _, guess := range guesses {
		if err := ValidateGuess(guess); err != nil {
			return err
		}
	}
	return nil
}

// TODO: Implement actual team validation logic.
func ValidateTeamId(teamid string) error {
	// Placeholder: always returns nil (valid)
	// In production, this would check against a database or auth service
	if teamid == "" {
		return fmt.Errorf("%w: team_id cannot be empty", ErrInvalidTeamId)
	}
	return nil
}
