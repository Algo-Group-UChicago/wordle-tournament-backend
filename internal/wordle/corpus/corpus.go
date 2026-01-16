package corpus

import (
	_ "embed"
	"encoding/csv"
	"log"
	"math"
	"strconv"
	"strings"
	"sync"
)

type wordSet map[string]struct{}
type answerWeightsMap map[string]float64

//go:embed corpus.txt
var corpusData string

//go:embed possible_answers.csv
var answersCsvData string

var (
	corpus          wordSet
	possibleAnswers answerWeightsMap
	once            sync.Once
)

func GetCorpus() wordSet {
	once.Do(initializeCorpus)
	return corpus
}

func GetGradingAnswerKey() answerWeightsMap {
	once.Do(initializeCorpus)
	return possibleAnswers
}

// GetWordWeight returns the weight for a given word.
// Returns 0.0 if the word is not found in the weights map.
func GetWordWeight(word string) float64 {
	once.Do(initializeCorpus)
	if weight, exists := possibleAnswers[word]; exists {
		return weight
	}
	return 0.0
}

func IsValidWord(word string) bool {
	_, exists := GetCorpus()[word]
	return exists
}

func initializeCorpus() {
	corpus = loadToSet(corpusData)
	possibleAnswers = loadWeightedAnswers(answersCsvData)
	log.Printf("Loaded %d words from corpus and %d possible answers with weights", len(corpus), len(possibleAnswers))
}

func loadToSet(data string) wordSet {
	ws := make(wordSet)
	for _, word := range strings.Fields(strings.ReplaceAll(data, "\n", " ")) {
		ws[word] = struct{}{}
	}
	return ws
}

func loadWeightedAnswers(csvData string) answerWeightsMap {
	reader := csv.NewReader(strings.NewReader(csvData))
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Failed to parse CSV: %v", err)
	}

	if len(records) == 0 {
		log.Fatal("CSV file is empty")
	}

	// Skip header row
	records = records[1:]

	// First pass: calculate log counts and find max
	result := make(answerWeightsMap, len(records))
	maxLogCount := 0.0

	for _, record := range records {
		if len(record) < 2 {
			continue
		}
		word := record[0]
		countStr := record[1]

		count, err := strconv.ParseInt(countStr, 10, 64)
		if err != nil {
			log.Printf("Warning: failed to parse count for word %s: %v", word, err)
			continue
		}

		logCount := math.Log(float64(count))
		result[word] = logCount
		if logCount > maxLogCount {
			maxLogCount = logCount
		}
	}

	// Second pass: normalize weights in place
	for word := range result {
		result[word] = result[word] / maxLogCount
	}

	return result
}
