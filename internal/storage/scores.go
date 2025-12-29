package storage

import "time"

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
	BestScore     float64        `dynamodbav:"best_score"`
	CompletedRuns []CompletedRun `dynamodbav:"completed_runs"`
}

// GetScore retrieves a ScoreItem from the Scores table by team_id.
// Returns a pointer to the item and nil error if found, or nil and an error if not found.
func GetScore(teamID string) (*ScoreItem, error)

// PutScore writes a ScoreItem to the Scores table in DynamoDB.
// Uses PutItem which will overwrite the entire item if it already exists, or create it if it doesn't.
// Returns an error if marshaling or writing to DynamoDB fails.
func PutScore(score *ScoreItem) error

// RemoveScore deletes a ScoreItem from the Scores table by team_id.
// Returns an error if the key marshaling or DeleteItem operation fails.
func RemoveScore(teamID string) error
