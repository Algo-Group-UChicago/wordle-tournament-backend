package common

import (
	"strconv"
	"time"

	"wordle-tournament-backend/internal/config"
)

// UnsolvedScoreSentinel is used to represent an unsolved run's score in DynamoDB
// (since DynamoDB cannot store infinity values)
const UnsolvedScoreSentinel = -1.0

// GetSeed returns the random seed to use for game generation.
// If RANDOM_SEED environment variable is set, uses that value.
// Otherwise, uses the current time in nanoseconds for randomness.
func GetSeed() int64 {
	cfg := config.Get()
	if cfg.RandomSeed != "" {
		seed, err := strconv.ParseInt(cfg.RandomSeed, 10, 64)
		if err != nil {
			return time.Now().UnixNano()
		}
		return seed
	}
	return time.Now().UnixNano()
}
