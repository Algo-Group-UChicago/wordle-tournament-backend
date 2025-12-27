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

// PutDefaultActiveRun creates a new ActiveRuns entry in DynamoDB for the given
// team_id and run_id. The entry contains a list of GameState entries, each with
// a unique randomly selected answer from the corpus. The item is configured with
// a TTL that expires after ActiveRunTTL duration.
//
// Returns an error if marshaling or writing to DynamoDB fails.
func PutDefaultActiveRun(teamID, runID string) error {
	item := ActiveRunItem{
		TeamID: teamID,
		RunID:  runID,
		Games:  createDefaultGameStates(),
		TTL:    time.Now().Add(ActiveRunTTL).Unix(),
	}

	return PutActiveRun(&item)
}

// GetActiveRun queries the ActiveRuns table by team_id and run_id to retrieve
// an ActiveRunItem. If the item is found, returns a pointer to the item and nil error.
// If the item is not found in the database, returns a nil pointer and an error.
func GetActiveRun(teamID, runID string) (*ActiveRunItem, error) {
	ctx := context.Background()

	client := getDynamoClient()

	key, err := attributevalue.MarshalMap(map[string]string{
		"team_id": teamID,
		"run_id":  runID,
	})
	if err != nil {
		return nil, fmt.Errorf("marshal key: %w", err)
	}

	result, err := client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(activeRunsTableName),
		Key:       key,
	})
	if err != nil {
		return nil, fmt.Errorf("DynamoDB GetItem operation failed: %w", err)
	}

	if result.Item == nil {
		return nil, fmt.Errorf("ActiveRuns item not found for team_id=%s, run_id=%s", teamID, runID)
	}

	var item ActiveRunItem
	if err := attributevalue.UnmarshalMap(result.Item, &item); err != nil {
		return nil, fmt.Errorf("unmarshal ActiveRuns item: %w", err)
	}

	return &item, nil
}

// PutActiveRun writes the provided ActiveRunItem to the ActiveRuns table in DynamoDB.
// It uses PutItem which will overwrite the entire item if it already exists, or create
// it if it doesn't. Returns an error if marshaling or writing to DynamoDB fails.
func PutActiveRun(activeRun *ActiveRunItem) error {
	ctx := context.Background()

	client := getDynamoClient()

	av, err := attributevalue.MarshalMap(activeRun)
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

// RemoveActiveRun deletes an ActiveRuns item from DynamoDB by team_id and run_id.
// Returns an error if the key marshaling or DeleteItem operation fails.
func RemoveActiveRun(teamID, runID string) error {
	ctx := context.Background()

	client := getDynamoClient()

	key, err := attributevalue.MarshalMap(map[string]string{
		"team_id": teamID,
		"run_id":  runID,
	})
	if err != nil {
		return fmt.Errorf("marshal key: %w", err)
	}

	_, err = client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(activeRunsTableName),
		Key:       key,
	})
	if err != nil {
		return fmt.Errorf("DynamoDB DeleteItem operation failed: %w", err)
	}

	return nil
}
