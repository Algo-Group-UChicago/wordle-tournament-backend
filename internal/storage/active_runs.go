package storage

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"wordle-tournament-backend/internal/common"
	"wordle-tournament-backend/internal/wordle/corpus"
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

func createDefaultRounds() []RoundState {
	possibleAnswers := corpus.GetGradingAnswerKey()

	// rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	rng := rand.New(rand.NewSource(1))
	selected := make(map[string]bool)
	rounds := make([]RoundState, 0, common.NumTargetWords)

	for len(rounds) < common.NumTargetWords {
		idx := rng.Intn(len(possibleAnswers))
		answer := possibleAnswers[idx]

		if !selected[answer] {
			selected[answer] = true
			rounds = append(rounds, RoundState{
				Solved:     false,
				NumGuesses: 0,
				Answer:     answer,
			})
		}
	}

	return rounds
}

func CreateDefaultActiveRun(teamID, runID string) error {
	ctx := context.Background()

	client := getDynamoClient()

	item := ActiveRunItem{
		TeamID: teamID,
		RunID:  runID,
		Rounds: createDefaultRounds(),
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
