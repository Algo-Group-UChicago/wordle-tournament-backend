package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

const (
	scoresTableName = "Scores"
)

// CompletedRun represents a single completed run for a team.
type CompletedRun struct {
	RunID          string    `dynamodbav:"run_id"`
	TotalScore     float64   `dynamodbav:"total_score"`
	AverageGuesses float64   `dynamodbav:"average_guesses"`
	SolvedCount    int       `dynamodbav:"solved_count"`
	CompletedAt    time.Time `dynamodbav:"completed_at"`
}

// ScoreItem maps a team_id to a list of CompletedRun entries in the Scores table.
// The table uses team_id as the partition key.
type ScoreItem struct {
	TeamID        string         `dynamodbav:"team_id"`
	CompletedRuns []CompletedRun `dynamodbav:"completed_runs"`
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
		return nil, fmt.Errorf("score not found for team_id=%s", teamID)
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
