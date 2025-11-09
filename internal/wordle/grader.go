package wordle

const WordLength = 5

// GradeGuess evaluates a guess against the answer and returns a hint string.
// The hint string is comma-separated with values: "correct", "present", or "absent".
// Uses the same logic as the Rust implementation with a two-pass algorithm:
// 1. First pass: mark all exact matches as "correct"
// 2. Second pass: mark remaining characters as "present" if they exist in unseen pool
func GradeGuess(guess, answer string) string {
	// Validate word lengths
	if len(guess) != WordLength || len(answer) != WordLength {
		panic("guess and answer must be exactly 5 characters")
	}

	// Initialize hint array with "absent"
	hintArr := [WordLength]string{"absent", "absent", "absent", "absent", "absent"}
	unseenPool := []rune{}

	// Convert strings to rune slices for proper character handling
	guessRunes := []rune(guess)
	answerRunes := []rune(answer)

	// Mark greens (correct position)
	for i := 0; i < WordLength; i++ {
		if guessRunes[i] == answerRunes[i] {
			hintArr[i] = "correct"
		} else {
			unseenPool = append(unseenPool, answerRunes[i])
		}
	}

	// Mark yellows (present but wrong position)
	for i := 0; i < WordLength; i++ {
		if hintArr[i] == "absent" {
			// Check if this character exists in unseen pool
			for j, unseenChar := range unseenPool {
				if guessRunes[i] == unseenChar {
					// Remove from unseen pool
					unseenPool = append(unseenPool[:j], unseenPool[j+1:]...)
					hintArr[i] = "present"
					break
				}
			}
		}
	}

	// Convert hint array to comma-separated string
	result := ""
	for i, hint := range hintArr {
		if i > 0 {
			result += ","
		}
		result += hint
	}

	return result
}
