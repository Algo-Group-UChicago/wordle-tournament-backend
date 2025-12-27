package common

import (
	"strconv"
	"time"

	"wordle-tournament-backend/internal/config"
)

// GetSeed returns the random seed to use for game generation.
// If RANDOM_SEED environment variable is set, uses that value.
// Otherwise, uses the current time in nanoseconds for randomness.
func GetSeed() int64 {
	cfg := config.Get()
	if cfg.RandomSeed != "" {
		seed, err := strconv.ParseInt(cfg.RandomSeed, 10, 64)
		if err != nil {
			// If seed is invalid, fall back to time-based seed
			return time.Now().UnixNano()
		}
		return seed
	}
	return time.Now().UnixNano()
}
