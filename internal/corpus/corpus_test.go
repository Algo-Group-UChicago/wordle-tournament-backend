package corpus

import (
	"testing"

	"wordle-tournament-backend/internal/common"
)

func TestCorpusLoads(t *testing.T) {
	corpus := GetCorpus()

	if len(corpus) == 0 {
		t.Error("Corpus should not be empty")
	}

	if len(corpus) <= 1000 {
		t.Errorf("Corpus should have many words, got %d", len(corpus))
	}
}

func TestAllWordsCorrectLength(t *testing.T) {
	corpus := GetCorpus()

	for word := range corpus {
		if len(word) != common.WordLength {
			t.Errorf("Word '%s' is not %d letters", word, common.WordLength)
		}
	}
}

func TestIsValidWord(t *testing.T) {
	tests := []struct {
		word  string
		valid bool
	}{
		{"crane", true},
		{"hello", true},
		{"zzzzz", false},
		{"notinlist", false},
	}

	for _, tt := range tests {
		t.Run(tt.word, func(t *testing.T) {
			result := IsValidWord(tt.word)
			if result != tt.valid {
				t.Errorf("IsValidWord(%q) = %v, want %v", tt.word, result, tt.valid)
			}
		})
	}
}

func TestCaseSensitive(t *testing.T) {
	// Assuming corpus is lowercase
	if !IsValidWord("crane") {
		t.Error("Expected 'crane' to be valid")
	}

	if IsValidWord("CRANE") {
		t.Error("Expected 'CRANE' to be invalid (case sensitive)")
	}
}
