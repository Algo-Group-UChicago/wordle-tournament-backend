package corpus

import (
	_ "embed"
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"sync"

	"wordle-tournament-backend/internal/common"
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
	common.LogInfo("corpus", fmt.Sprintf("Loaded %d words from corpus and %d possible answers with weights", len(corpus), len(possibleAnswers)))
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
		common.LogError("corpus", "Failed to parse CSV", 0, slog.String("error", err.Error()))
		os.Exit(1)
	}

	if len(records) == 0 {
		common.LogError("corpus", "CSV file is empty", 0)
		os.Exit(1)
	}

	// Skip header row
	records = records[1:]

	result := make(answerWeightsMap, len(records))

	for _, record := range records {
		if len(record) < 3 {
			continue
		}
		word := record[0]
		weightStr := record[2]

		weight, err := strconv.ParseFloat(weightStr, 64)
		if err != nil {
			common.LogWarning("corpus", "failed to parse weight for word", 0, slog.String("word", word), slog.String("error", err.Error()))
			continue
		}

		result[word] = weight
	}

	return result
}
