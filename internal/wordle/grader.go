package wordle

import (
	"strings"

	"wordle-tournament-backend/internal/common"
)

// GradeGuesses takes two lists of guesses and answers and returns an array of hints.
// Each hint is a 5-character string where:
// - 'O' indicates a correct letter in the correct position
// - '~' indicates a correct letter in the wrong position
// - 'X' indicates a letter not in the answer
func GradeGuesses(guesses, answers []string) []string {

	// TODO: Remove panic before deploying to production
	// Confirm guesses and answers are the same length
	if len(guesses) != len(answers) {
		panic("guesses and answers must be the same length")
	}

	hints := make([]string, len(guesses))
	for i := 0; i < len(guesses); i++ {
		// if we receive a DummyGuess from middleware
		// then the guess is automatically correct
		if guesses[i] == common.DummyGuess {
			hints[i] = strings.Repeat("O", common.WordLength)
		} else {
			hints[i] = gradeGuessLogical(guesses[i], answers[i])
		}
	}

	return hints
}

// grade a single guess and answer mirroring the rust algorithm
func gradeGuessLogical(guess, answer string) string {
	hint := []rune("XXXXX")
	remainingChars := []rune(answer)

	guessRunes := []rune(guess)
	answerRunes := []rune(answer)

	// Mark correctly placed characters and remove them from remainingChars
	for i := 0; i < common.WordLength; i++ {
		if guessRunes[i] == answerRunes[i] {
			hint[i] = 'O'
			removeFirstOccurrence(&remainingChars, guessRunes[i])
		}
	}

	// Go through remaining characters in guess
	// If it exists in the list, mark "~" and remove from list, otherwise leave as "X"
	for i := 0; i < common.WordLength; i++ {
		if hint[i] == 'X' { // Only check characters not already marked as correct
			if removeFirstOccurrence(&remainingChars, guessRunes[i]) {
				hint[i] = '~'
			}
		}
	}

	return string(hint)
}

func removeFirstOccurrence(slice *[]rune, target rune) bool {
	s := *slice
	for i, char := range s {
		if char == target {
			*slice = append(s[:i], s[i+1:]...)
			return true
		}
	}
	return false
}
