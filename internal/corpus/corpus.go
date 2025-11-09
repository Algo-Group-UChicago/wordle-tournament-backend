package corpus

import (
	_ "embed"
	"log"
	"strings"
	"sync"
)

//go:embed corpus.txt
var corpusData string

//go:embed possible_answers.txt
var answersData string

var (
	corpus          map[string]struct{}
	possibleAnswers []string
	once            sync.Once
)

func GetCorpus() map[string]struct{} {
	once.Do(initializeMaps)
	return corpus
}

func GetGradingAnswerKey() []string {
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
	possibleAnswers = loadToSlice(answersData)
	log.Printf("Loaded %d words from corpus and %d possible answers", len(corpus), len(possibleAnswers))
}

func loadFromString(data string) map[string]struct{} {
	wordSet := make(map[string]struct{})

	for _, line := range strings.Split(data, "\n") {
		if word := strings.TrimSpace(line); word != "" {
			wordSet[word] = struct{}{}
		}
	}

	return wordSet
}

func loadToSlice(data string) []string {
	var words []string

	lines := strings.Split(data, "\n")
	for _, word := range lines {
		word = strings.TrimSpace(word) // Remove any whitespace
		if word != "" {                // Skip empty lines
			words = append(words, word)
		}
	}

	return words
}
