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
	if len(guesses) != len(answers) {
		panic("guesses and answers must be the same length")
	}

	hints := make([]string, len(guesses))
	for i := 0; i < len(guesses); i++ {
		if guesses[i] == common.DummyGuess {
			hints[i] = strings.Repeat("O", common.WordLength)
		} else {
			hints[i] = gradeGuessLogical(guesses[i], answers[i])
		}
	}

	return hints
}

// Grade a single guess and answer mirroring the rust algorithm
func gradeGuessLogical(guess, answer string) string {
	hint := []rune("XXXXX")
	remainingChars := []rune(answer)

	guessRunes := []rune(guess)
	answerRunes := []rune(answer)

	// Mark correctly placed characters
	for i := 0; i < common.WordLength; i++ {
		if guessRunes[i] == answerRunes[i] {
			hint[i] = 'O'
			tryRemoveFirst(&remainingChars, guessRunes[i])
		}
	}

	// Mark misplaced characters
	for i := 0; i < common.WordLength; i++ {
		if hint[i] == 'X' {
			if tryRemoveFirst(&remainingChars, guessRunes[i]) {
				hint[i] = '~'
			}
		}
	}

	return string(hint)
}

func tryRemoveFirst(slice *[]rune, target rune) bool {
	s := *slice
	for i, char := range s {
		if char == target {
			*slice = append(s[:i], s[i+1:]...)
			return true
		}
	}
	return false
}
