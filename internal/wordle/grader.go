package wordle

import "strings"

const WordLength = 5
const DummyGuess = "imagine guessing more than 5 letters"

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
		if guesses[i] == DummyGuess {
			hints[i] = strings.Repeat("O", WordLength)
		} else {
			hints[i] = gradeGuessList(guesses[i], answers[i])
		}
	}

	return hints
}

// grade a single guess and answer mirroring the rust algorithm
// imo this impl is cumbersome to reason about
func gradeGuessList(guess, answer string) string {
	hint := []rune("XXXXX")
	remainingChars := []rune(answer)

	// Mark corrrectly placed characters and remove them from remainingChars
	guessRunes := []rune(guess)
	answerRunes := []rune(answer)
	for i := 0; i < WordLength; i++ {
		if guessRunes[i] == answerRunes[i] {
			hint[i] = 'O'
			// Remove this character from remainingChars
			for j, char := range remainingChars {
				if char == guessRunes[i] {
					remainingChars = append(remainingChars[:j], remainingChars[j+1:]...)
					break
				}
			}
		}
	}

	// Go through remaining characters in guess
	// If it exists in the list, mark "~" and remove from list, otherwise leave as "X"
	for i := 0; i < WordLength; i++ {
		if hint[i] == 'X' { // Only check characters not already marked as correct
			for j, char := range remainingChars {
				if guessRunes[i] == char {
					hint[i] = '~'
					remainingChars = append(remainingChars[:j], remainingChars[j+1:]...)
					break
				}
			}
		}
	}

	return string(hint)
}
