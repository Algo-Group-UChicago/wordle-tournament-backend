package corpus

import (
	_ "embed"
	"log"
	"strings"
	"sync"
)

type wordSet map[string]struct{}

//go:embed corpus.txt
var corpusData string

//go:embed possible_answers.txt
var answersData string

var (
	corpus          wordSet
	possibleAnswers []string
	once            sync.Once
)

func GetCorpus() wordSet {
	once.Do(initializeCorpus)
	return corpus
}

func GetGradingAnswerKey() []string {
	once.Do(initializeCorpus)
	return possibleAnswers
}

func IsValidWord(word string) bool {
	_, exists := GetCorpus()[word]
	return exists
}

func initializeCorpus() {
	corpus = loadToSet(corpusData)
	possibleAnswers = loadToSlice(answersData)
	log.Printf("Loaded %d words from corpus and %d possible answers", len(corpus), len(possibleAnswers))
}

func loadToSet(data string) wordSet {
	ws := make(wordSet)
	for _, word := range strings.Fields(strings.ReplaceAll(data, "\n", " ")) {
		ws[word] = struct{}{}
	}
	return ws
}

func loadToSlice(data string) []string {
	return strings.Fields(strings.ReplaceAll(data, "\n", " "))
}
