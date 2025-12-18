package storage

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

const activeRunsTableName = "ActiveRuns"

// RoundState represents a single Wordle game within a run.
type RoundState struct {
	Solved     bool   `json:"solved" dynamodbav:"solved"`
	NumGuesses int    `json:"num_guesses" dynamodbav:"num_guesses"`
	Answer     string `json:"answer" dynamodbav:"answer"`
}

type ActiveRunItem struct {
	TeamID string       `dynamodbav:"team_id"`
	RunID  string       `dynamodbav:"run_id"`
	Rounds []RoundState `dynamodbav:"rounds"`
}

// CreateBlankActiveRun creates an ActiveRuns row with an empty rounds list.
func CreateBlankActiveRun(ctx context.Context, teamID, runID string) error {
	if err != nil {
		return err
	}

	blank_item := ActiveRunItem{
		TeamID: teamID,
		RunID:  runID,
		Rounds: []RoundState{},
	}

	av, err := attributevalue.MarshalMap(blank_item)
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
