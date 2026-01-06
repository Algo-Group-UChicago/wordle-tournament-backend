package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

const (
	scoresTableName = "Scores"
)

var (
	// ErrScoreNotFound is returned when a score for a team_id is not found in the Scores table.
	ErrScoreNotFound = errors.New("score not found")
)

// CompletedRun represents a single completed run for a team.
type CompletedRun struct {
	RunID          string    `dynamodbav:"run_id"`
	Score          float64   `dynamodbav:"score"`
	AverageGuesses float64   `dynamodbav:"average_guesses"`
	Solved         bool      `dynamodbav:"solved"`
	CompletedAt    time.Time `dynamodbav:"completed_at"`
}

// ScoreItem maps a team_id to a list of CompletedRun entries in the Scores table.
// The table uses team_id as the partition key.
type ScoreItem struct {
	TeamID        string         `dynamodbav:"team_id"`
	CompletedRuns []CompletedRun `dynamodbav:"completed_runs"`
}

// CalculateScore calculates the total score from a list of GameState entries.
// Returns the average number of guesses across all games.
func CalculateScore(games []GameState) float64 {
	if len(games) == 0 {
		return 0.0
	}

	totalGuesses := 0.0
	for _, game := range games {
		totalGuesses += float64(game.NumGuesses)
	}

	return totalGuesses / float64(len(games))
}

// GetScore retrieves a ScoreItem from the Scores table by team_id.
// Returns a pointer to the item and nil error if found, or nil and an error if not found.
func GetScore(teamID string) (*ScoreItem, error) {
	ctx := context.Background()

	client := getDynamoClient()

	key, err := attributevalue.MarshalMap(map[string]string{
		"team_id": teamID,
	})
	if err != nil {
		return nil, fmt.Errorf("marshal key: %w", err)
	}

	result, err := client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(scoresTableName),
		Key:       key,
	})
	if err != nil {
		return nil, fmt.Errorf("DynamoDB GetItem operation failed: %w", err)
	}

	if result.Item == nil {
		return nil, fmt.Errorf("%w for team_id=%s", ErrScoreNotFound, teamID)
	}

	var item ScoreItem
	if err := attributevalue.UnmarshalMap(result.Item, &item); err != nil {
		return nil, fmt.Errorf("unmarshal Scores item: %w", err)
	}

	return &item, nil
}

// PutScore writes a ScoreItem to the Scores table in DynamoDB.
// Uses PutItem which will overwrite the entire item if it already exists, or create it if it doesn't.
// Returns an error if marshaling or writing to DynamoDB fails.
func PutScore(score *ScoreItem) error {
	ctx := context.Background()

	client := getDynamoClient()

	av, err := attributevalue.MarshalMap(score)
	if err != nil {
		return fmt.Errorf("marshal Scores item: %w", err)
	}

	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(scoresTableName),
		Item:      av,
	})
	if err != nil {
		return fmt.Errorf("put Scores item: %w", err)
	}

	return nil
}
