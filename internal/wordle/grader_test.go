package wordle

import (
	"testing"
)

func TestAllAbsent(t *testing.T) {
	result := GradeGuesses([]string{"crane"}, []string{"built"})
	expected := "XXXXX"
	if result[0] != expected {
		t.Errorf("expected %q, got %q", expected, result[0])
	}
}

func TestAllCorrect(t *testing.T) {
	result := GradeGuesses([]string{"crane"}, []string{"crane"})
	expected := "OOOOO"
	if result[0] != expected {
		t.Errorf("expected %q, got %q", expected, result[0])
	}
}

func TestDuplicateInTargetCausesPresent(t *testing.T) {
	result := GradeGuesses([]string{"roost"}, []string{"robot"})
	expected := "OO~XO"
	if result[0] != expected {
		t.Errorf("expected %q, got %q", expected, result[0])
	}
}

func TestGuessHasMoreDuplicatesThanTarget(t *testing.T) {
	result := GradeGuesses([]string{"allee"}, []string{"apple"})
	expected := "O~XXO"
	if result[0] != expected {
		t.Errorf("expected %q, got %q", expected, result[0])
	}
}

func TestNoMatches(t *testing.T) {
	result := GradeGuesses([]string{"crane"}, []string{"yummy"})
	expected := "XXXXX"
	if result[0] != expected {
		t.Errorf("expected %q, got %q", expected, result[0])
	}
}

func TestMixedDuplicatesAndCorrect(t *testing.T) {
	result := GradeGuesses([]string{"ABBEY"}, []string{"BANAL"})
	expected := "~~XXX"
	if result[0] != expected {
		t.Errorf("expected %q, got %q", expected, result[0])
	}
}

func TestDuplicateCorrectAndPresentSameLetter(t *testing.T) {
	result := GradeGuesses([]string{"array"}, []string{"alarm"})
	expected := "O~X~X"
	if result[0] != expected {
		t.Errorf("expected %q, got %q", expected, result[0])
	}
}

func TestPresentDoesNotStealFromCorrect(t *testing.T) {
	result := GradeGuesses([]string{"babee"}, []string{"aback"})
	expected := "~~XXX"
	if result[0] != expected {
		t.Errorf("expected %q, got %q", expected, result[0])
	}
}
