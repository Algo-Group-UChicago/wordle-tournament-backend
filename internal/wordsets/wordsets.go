package wordsets

import (
	"bufio"
	"log"
	"os"
	"sync"
)

// Package-level variables for the word sets
var (
	corpus          map[string]struct{}
	possibleAnswers map[string]struct{}
	once            sync.Once
)

// Initialize loads the word sets from files exactly once
// This is thread-safe and will only execute once even if called concurrently
func Initialize(corpusPath, answersPath string) {
	once.Do(func() {
		corpus = loadWordSet(corpusPath)
		possibleAnswers = loadWordSet(answersPath)
		log.Printf("Loaded %d words from corpus and %d possible answers", len(corpus), len(possibleAnswers))
	})
}

// IsInCorpus checks if a word exists in the corpus
// This is safe for concurrent access (read-only after initialization)
func IsInCorpus(word string) bool {
	_, exists := corpus[word]
	return exists
}

// IsValidAnswer checks if a word is a valid Wordle answer
// This is safe for concurrent access (read-only after initialization)
func IsValidAnswer(word string) bool {
	_, exists := possibleAnswers[word]
	return exists
}

// loadWordSet reads a file and returns a set (map[string]struct{}) of words
func loadWordSet(filepath string) map[string]struct{} {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("Failed to open word file %s: %v", filepath, err)
	}
	defer file.Close()

	wordSet := make(map[string]struct{})
	scanner := bufio.Scanner(file)

	for scanner.Scan() {
		word := scanner.Text()
		wordSet[word] = struct{}{} // empty struct takes 0 bytes
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading word file %s: %v", filepath, err)
	}

	return wordSet
}
