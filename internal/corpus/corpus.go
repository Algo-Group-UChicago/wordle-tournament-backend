package corpus

import (
	"bufio"
	"log"
	"os"
	"sync"
)

const CorpusFile = "corpus.txt"
const AnswerKeyFile = "possible_answers.txt"

var (
	corpus          map[string]struct{}
	possibleAnswers map[string]struct{}
	once            sync.Once
)

func GetCorpus() map[string]struct{} {
	once.Do(func() {
		corpus = loadDictionary(CorpusFile)
		possibleAnswers = loadDictionary(AnswerKeyFile)
		log.Printf("Loaded %d words from corpus and %d possible answers", len(corpus), len(possibleAnswers))
	})
	return corpus
}

func GetGradingAnswerKey() map[string]struct{} {
	once.Do(func() {
		corpus = loadDictionary("corpus.txt")
		possibleAnswers = loadDictionary("possible_answers.txt")
		log.Printf("Loaded %d words from corpus and %d possible answers", len(corpus), len(possibleAnswers))
	})
	return possibleAnswers
}

func IsValidWord(word string) bool {
	_, exists := GetCorpus()[word]
	return exists
}

func loadDictionary(filepath string) map[string]struct{} {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("Failed to open word file %s: %v", filepath, err)
	}
	defer file.Close()

	wordSet := make(map[string]struct{})
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		word := scanner.Text()
		wordSet[word] = struct{}{} // empty struct takes 0 bytes
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading word file %s: %v", filepath, err)
	}

	return wordSet
}
