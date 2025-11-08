package registry

import (
	"bufio"
	"log"
	"os"
	"sync"
)

var (
	corpus          map[string]struct{}
	possibleAnswers map[string]struct{}
	once            sync.Once
)

func Initialize(corpusPath, answersPath string) {
	once.Do(func() {
		corpus = loadDictionary(corpusPath)
		possibleAnswers = loadDictionary(answersPath)
		log.Printf("Loaded %d words from corpus and %d possible answers", len(corpus), len(possibleAnswers))
	})
}

func IsInCorpus(word string) bool {
	_, exists := corpus[word]
	return exists
}

func IsValidAnswer(word string) bool {
	_, exists := possibleAnswers[word]
	return exists
}

func loadDictionary(filepath string) map[string]struct{} {
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
