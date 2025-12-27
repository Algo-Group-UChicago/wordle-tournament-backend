package storage

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"wordle-tournament-backend/internal/common"
	"wordle-tournament-backend/internal/wordle/corpus"
)

const (
	activeRunsTableName = "ActiveRuns"
	ActiveRunTTL        = 10 * time.Minute
)

// GameState represents a single Wordle game within a run.
type GameState struct {
	Solved     bool   `json:"solved" dynamodbav:"solved"`
	NumGuesses int    `json:"num_guesses" dynamodbav:"num_guesses"`
	Answer     string `json:"answer" dynamodbav:"answer"`
}

// ActiveRunItem maps (team_id, run_id) to a list of GameState entries with TTL.
type ActiveRunItem struct {
	TeamID string      `dynamodbav:"team_id"`
	RunID  string      `dynamodbav:"run_id"`
	Games  []GameState `dynamodbav:"games"`
	TTL    int64       `dynamodbav:"ttl"`
}

// createDefaultGames returns a slice of GameState entries, each containing
// a unique randomly selected answer from the corpus.
func createDefaultGameStates() []GameState {
	possibleAnswers := corpus.GetGradingAnswerKey()

	// rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	// Creates identical games for each run
	rng := rand.New(rand.NewSource(1))
	selected := make(map[string]bool)
	games := make([]GameState, 0, common.NumTargetWords)

	for len(games) < common.NumTargetWords {
		idx := rng.Intn(len(possibleAnswers))
		answer := possibleAnswers[idx]

		if !selected[answer] {
			selected[answer] = true
			games = append(games, GameState{
				Solved:     false,
				NumGuesses: 0,
				Answer:     answer,
			})
		}
	}

	return games
}

// CreateDefaultActiveRun creates a new ActiveRuns entry in DynamoDB for the given
// team_id and run_id. The entry contains a list of GameState entries, each with
// a unique randomly selected answer from the corpus. The item is configured with
// a TTL that expires after ActiveRunTTL duration.
//
// Returns an error if marshaling or writing to DynamoDB fails.
func CreateDefaultActiveRun(teamID, runID string) error {
	ctx := context.Background()

	client := getDynamoClient()

	item := ActiveRunItem{
		TeamID: teamID,
		RunID:  runID,
		Games:  createDefaultGameStates(),
		TTL:    time.Now().Add(ActiveRunTTL).Unix(),
	}

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("marshal ActiveRuns item: %w", err)
	}

	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(activeRunsTableName),
		Item:      av,
	})
	if err != nil {
		return fmt.Errorf("put ActiveRuns item: %w", err)
	}

	return nil
}
