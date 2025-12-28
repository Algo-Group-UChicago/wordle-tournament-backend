//go:build integration

package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"wordle-tournament-backend/internal/common"
	"wordle-tournament-backend/internal/handlers"
	"wordle-tournament-backend/internal/server"
	"wordle-tournament-backend/internal/storage"
)

func setupIntegrationTest(t *testing.T) *httptest.Server {
	srv := server.New()
	ts := httptest.NewServer(srv.Handler())
	return ts
}

// TestIntegrationInstantSolve tests the complete flow of starting a run and solving all games
// instantly with perfect guesses. It performs the following steps:
//  1. Starts a new run by calling POST /start with a team_id, receiving a run_id
//  2. Verifies the run was created with the correct number of games (NumTargetWords),
//     and that all games are initially unsolved with 0 guesses and valid answers
//  3. Builds a guesses array using the actual answers from each game
//  4. Submits all perfect guesses via POST /api/guesses
//  5. Verifies the response contains hints for all games, and all hints are "OOOOO" (all correct)
//  6. Verifies all games are now marked as solved with NumGuesses = 1
func TestIntegrationInstantSolve(t *testing.T) {
	ts := setupIntegrationTest(t)
	defer ts.Close()

	teamID := "TEST_TEAM"
	client := &http.Client{}

	// Step 1: Start a new run
	startReq := handlers.StartRequest{TeamID: teamID}
	startBody, err := json.Marshal(startReq)
	if err != nil {
		t.Fatalf("Failed to marshal start request: %v", err)
	}

	startResp, err := client.Post(ts.URL+"/start", "application/json", bytes.NewBuffer(startBody))
	if err != nil {
		t.Fatalf("Failed to call /start: %v", err)
	}
	defer startResp.Body.Close()

	if startResp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status 201 Created, got %d", startResp.StatusCode)
	}

	var startResponse handlers.StartResponse
	if err := json.NewDecoder(startResp.Body).Decode(&startResponse); err != nil {
		t.Fatalf("Failed to decode start response: %v", err)
	}

	if startResponse.RunID == "" {
		t.Fatal("RunID should not be empty")
	}

	runID := startResponse.RunID
	t.Logf("Created run with ID: %s", runID)

	// Step 2: Verify the run was created with correct number of games
	activeRun, err := storage.GetActiveRun(teamID, runID)
	if err != nil {
		t.Fatalf("Failed to get active run: %v", err)
	}

	if len(activeRun.Games) != common.NumTargetWords {
		t.Fatalf("Expected %d games, got %d", common.NumTargetWords, len(activeRun.Games))
	}

	// Verify all games are initially unsolved with 0 guesses and valid answers
	for i, game := range activeRun.Games {
		if game.Solved {
			t.Errorf("Game %d should not be solved initially", i)
		}
		if game.NumGuesses != 0 {
			t.Errorf("Game %d should have 0 guesses initially, got %d", i, game.NumGuesses)
		}
		if game.Answer == "" {
			t.Errorf("Game %d should have an answer", i)
		}
	}

	// Step 3: Build guesses array using the actual answers from each game
	guesses := make([]string, common.NumTargetWords)
	for i := range activeRun.Games {
		guesses[i] = activeRun.Games[i].Answer
	}

	// Step 4: Submit all perfect guesses via POST /api/guesses
	guessesReq := handlers.GuessesRequest{
		TeamId:  teamID,
		RunId:   runID,
		Guesses: guesses,
	}

	guessesBody, err := json.Marshal(guessesReq)
	if err != nil {
		t.Fatalf("Failed to marshal guesses request: %v", err)
	}

	guessesResp, err := client.Post(ts.URL+"/api/guesses", "application/json", bytes.NewBuffer(guessesBody))
	if err != nil {
		t.Fatalf("Failed to call /api/guesses: %v", err)
	}
	defer guessesResp.Body.Close()

	if guessesResp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200 OK, got %d", guessesResp.StatusCode)
	}

	var guessesResponse handlers.GuessesResponse
	if err := json.NewDecoder(guessesResp.Body).Decode(&guessesResponse); err != nil {
		t.Fatalf("Failed to decode guesses response: %v", err)
	}

	if len(guessesResponse.Hints) != common.NumTargetWords {
		t.Fatalf("Expected %d hints, got %d", common.NumTargetWords, len(guessesResponse.Hints))
	}

	// Step 5: Verify the response contains hints for all games, and all hints are "OOOOO" (all correct)
	for i, hint := range guessesResponse.Hints {
		if hint != "OOOOO" {
			t.Errorf("Game %d should have hint 'OOOOO' for perfect guess, got %q", i, hint)
		}
	}

	// Step 6: Verify all games are now marked as solved with NumGuesses = 1
	activeRunAfter, err := storage.GetActiveRun(teamID, runID)
	if err != nil {
		t.Fatalf("Failed to get active run after guesses: %v", err)
	}

	if len(activeRunAfter.Games) != common.NumTargetWords {
		t.Fatalf("Expected %d games after guesses, got %d", common.NumTargetWords, len(activeRunAfter.Games))
	}

	allSolved := true
	for i, game := range activeRunAfter.Games {
		if !game.Solved {
			t.Errorf("Game %d should be solved after submitting perfect guess", i)
			allSolved = false
		}
		if game.NumGuesses != 1 {
			t.Errorf("Game %d should have 1 guess after perfect guess, got %d", i, game.NumGuesses)
		}
	}

	if !allSolved {
		t.Fatal("Not all games were marked as solved")
	}

	t.Logf("Successfully processed all %d words", common.NumTargetWords)
}
