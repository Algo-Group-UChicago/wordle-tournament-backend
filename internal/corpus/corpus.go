package corpus

import (
	_ "embed"
	"log"
	"strings"
	"sync"
)

// Embed the corpus files directly into the binary at compile time
// Files are located in the project root, so we use ../../ to go up from internal/corpus/
//
//go:embed corpus.txt
var corpusData string

//go:embed possible_answers.txt
var answersData string

var (
	corpus          map[string]struct{}
	possibleAnswers map[string]struct{}
	once            sync.Once
)

func GetCorpus() map[string]struct{} {
	once.Do(initializeMaps)
	return corpus
}

func GetGradingAnswerKey() map[string]struct{} {
	once.Do(initializeMaps)
	return possibleAnswers
}

func IsValidWord(word string) bool {
	_, exists := GetCorpus()[word]
	return exists
}

// initializeMaps loads both corpus and answers from embedded data
func initializeMaps() {
	corpus = loadFromString(corpusData)
	possibleAnswers = loadFromString(answersData)
	log.Printf("Loaded %d words from corpus and %d possible answers", len(corpus), len(possibleAnswers))
}

// loadFromString parses embedded string data into a word set
func loadFromString(data string) map[string]struct{} {
	wordSet := make(map[string]struct{})

	// Split by newlines and add each word
	lines := strings.Split(data, "\n")
	for _, word := range lines {
		word = strings.TrimSpace(word) // Remove any whitespace
		if word != "" {                // Skip empty lines
			wordSet[word] = struct{}{}
		}
	}

	return wordSet
}
