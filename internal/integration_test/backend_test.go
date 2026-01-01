//go:build integration

package integration_test

//AI Generated Integration Tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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
//  7. Calls POST /end to complete the run
//  8. Verifies the correct entry is created in Scores database
//  9. Verifies the ActiveRun entry has been removed
//  10. Calls POST /end again with the same team_id and run_id, which should fail since the ActiveRun no longer exists
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

	// Step 7: Call POST /end to complete the run
	endReq := handlers.EndRequest{
		TeamID: teamID,
		RunID:  runID,
	}
	endBody, err := json.Marshal(endReq)
	if err != nil {
		t.Fatalf("Failed to marshal end request: %v", err)
	}

	endResp, err := client.Post(ts.URL+"/end", "application/json", bytes.NewBuffer(endBody))
	if err != nil {
		t.Fatalf("Failed to call /end: %v", err)
	}
	defer endResp.Body.Close()

	if endResp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200 OK, got %d", endResp.StatusCode)
	}

	var endResponse handlers.EndResponse
	if err := json.NewDecoder(endResp.Body).Decode(&endResponse); err != nil {
		t.Fatalf("Failed to decode end response: %v", err)
	}

	// Verify response values
	if !endResponse.Solved {
		t.Error("End response should indicate all games are solved")
	}
	if endResponse.AverageGuesses != 1.0 {
		t.Errorf("End response averageGuesses should be 1.0 (all games solved in 1 guess), got %f", endResponse.AverageGuesses)
	}
	if endResponse.Score <= 0 {
		t.Errorf("End response score should be positive, got %f", endResponse.Score)
	}
	if endResponse.Score != endResponse.AverageGuesses {
		t.Errorf("End response score and averageGuesses should be equal, got score=%f, averageGuesses=%f", endResponse.Score, endResponse.AverageGuesses)
	}

	// Step 8: Verify the correct entry is created in Scores database
	scoreItem, err := storage.GetScore(teamID)
	if err != nil {
		t.Fatalf("Failed to get score: %v", err)
	}

	if scoreItem.TeamID != teamID {
		t.Errorf("ScoreItem TeamID should be %s, got %s", teamID, scoreItem.TeamID)
	}

	if len(scoreItem.CompletedRuns) != 1 {
		t.Fatalf("Expected 1 completed run, got %d", len(scoreItem.CompletedRuns))
	}

	completedRun := scoreItem.CompletedRuns[0]
	if completedRun.RunID != runID {
		t.Errorf("CompletedRun RunID should be %s, got %s", runID, completedRun.RunID)
	}
	if !completedRun.Solved {
		t.Error("CompletedRun should be marked as solved")
	}
	if completedRun.Score != endResponse.Score {
		t.Errorf("CompletedRun Score should be %f, got %f", endResponse.Score, completedRun.Score)
	}
	if completedRun.AverageGuesses != endResponse.AverageGuesses {
		t.Errorf("CompletedRun AverageGuesses should be %f, got %f", endResponse.AverageGuesses, completedRun.AverageGuesses)
	}

	// Step 9: Verify the ActiveRun entry has been removed
	_, err = storage.GetActiveRun(teamID, runID)
	if err == nil {
		t.Fatal("ActiveRun should have been removed after calling /end")
	}
	if err != nil && !strings.Contains(err.Error(), "expired or not found") {
		t.Fatalf("Expected 'expired or not found' error, got: %v", err)
	}

	// Step 10: Call POST /end again with the same team_id and run_id, which should fail
	endReq2 := handlers.EndRequest{
		TeamID: teamID,
		RunID:  runID,
	}
	endBody2, err := json.Marshal(endReq2)
	if err != nil {
		t.Fatalf("Failed to marshal end request: %v", err)
	}

	endResp2, err := client.Post(ts.URL+"/end", "application/json", bytes.NewBuffer(endBody2))
	if err != nil {
		t.Fatalf("Failed to call /end: %v", err)
	}
	defer endResp2.Body.Close()

	// Verify /end fails with 400 Bad Request since ActiveRun no longer exists
	if endResp2.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected status 400 Bad Request when ActiveRun doesn't exist, got %d", endResp2.StatusCode)
	}

	t.Logf("Successfully processed all %d words", common.NumTargetWords)
}

// TestIntegrationMultipleGuessRounds tests the behavior of multiple rounds of guesses:
//  1. Starts a new run
//  2. Submits incorrect guesses for the first 3 games, verifying NumGuesses increments
//  3. Submits correct guesses for the first 2 games, verifying they become solved
//  4. Submits another round with correct guess for game 3, verifying it becomes solved
//  5. Verifies that already-solved games (games 0 and 1) don't increment NumGuesses
//     when submitting another round of guesses
//  6. Calls POST /end once at the end to complete the run (all games are now solved)
//  7. Verifies the correct entry is created in Scores database
//  8. Verifies the ActiveRun entry has been removed
func TestIntegrationMultipleGuessRounds(t *testing.T) {
	ts := setupIntegrationTest(t)
	defer ts.Close()

	teamID := "TEST_TEAM_MULTIPLE"
	client := &http.Client{}

	// Step 1: Start a new run
	startReq := handlers.StartRequest{TeamID: teamID}
	startBody, _ := json.Marshal(startReq)
	startResp, err := client.Post(ts.URL+"/start", "application/json", bytes.NewBuffer(startBody))
	if err != nil {
		t.Fatalf("Failed to call /start: %v", err)
	}
	defer startResp.Body.Close()

	var startResponse handlers.StartResponse
	json.NewDecoder(startResp.Body).Decode(&startResponse)
	runID := startResponse.RunID

	// Get the active run to access game answers
	activeRun, err := storage.GetActiveRun(teamID, runID)
	if err != nil {
		t.Fatalf("Failed to get active run: %v", err)
	}

	if len(activeRun.Games) < 3 {
		t.Fatalf("Need at least 3 games for this test, got %d", len(activeRun.Games))
	}

	// Step 2: Submit incorrect guesses for first 3 games
	guesses := make([]string, common.NumTargetWords)
	guesses[0] = "crane" // Wrong guess for game 0
	guesses[1] = "house" // Wrong guess for game 1
	guesses[2] = "table" // Wrong guess for game 2
	for i := 3; i < common.NumTargetWords; i++ {
		guesses[i] = common.DummyGuess
	}

	guessesReq := handlers.GuessesRequest{
		TeamId:  teamID,
		RunId:   runID,
		Guesses: guesses,
	}

	guessesBody, _ := json.Marshal(guessesReq)
	guessesResp, err := client.Post(ts.URL+"/api/guesses", "application/json", bytes.NewBuffer(guessesBody))
	if err != nil {
		t.Fatalf("Failed to call /api/guesses: %v", err)
	}
	defer guessesResp.Body.Close()

	// Verify NumGuesses incremented for first 3 games
	activeRunAfter1, err := storage.GetActiveRun(teamID, runID)
	if err != nil {
		t.Fatalf("Failed to get active run: %v", err)
	}

	for i := 0; i < 3; i++ {
		if activeRunAfter1.Games[i].NumGuesses != 1 {
			t.Errorf("Game %d should have 1 guess after first round, got %d", i, activeRunAfter1.Games[i].NumGuesses)
		}
		if activeRunAfter1.Games[i].Solved {
			t.Errorf("Game %d should not be solved after incorrect guess", i)
		}
	}

	// Step 3: Submit correct guesses for first 2 games
	guesses[0] = activeRun.Games[0].Answer // Correct for game 0
	guesses[1] = activeRun.Games[1].Answer // Correct for game 1
	guesses[2] = "wrong"                   // Still wrong for game 2
	// Games 3+ continue using DummyGuess (once DummyGuess is used, it's always used)

	guessesBody, _ = json.Marshal(guessesReq)
	guessesResp, err = client.Post(ts.URL+"/api/guesses", "application/json", bytes.NewBuffer(guessesBody))
	if err != nil {
		t.Fatalf("Failed to call /api/guesses: %v", err)
	}
	defer guessesResp.Body.Close()

	activeRunAfter2, err := storage.GetActiveRun(teamID, runID)
	if err != nil {
		t.Fatalf("Failed to get active run: %v", err)
	}

	// Verify games 0 and 1 are solved with NumGuesses = 2
	if !activeRunAfter2.Games[0].Solved {
		t.Error("Game 0 should be solved after correct guess")
	}
	if activeRunAfter2.Games[0].NumGuesses != 2 {
		t.Errorf("Game 0 should have 2 guesses, got %d", activeRunAfter2.Games[0].NumGuesses)
	}

	if !activeRunAfter2.Games[1].Solved {
		t.Error("Game 1 should be solved after correct guess")
	}
	if activeRunAfter2.Games[1].NumGuesses != 2 {
		t.Errorf("Game 1 should have 2 guesses, got %d", activeRunAfter2.Games[1].NumGuesses)
	}

	// Verify game 2 is still unsolved with NumGuesses = 2
	if activeRunAfter2.Games[2].Solved {
		t.Error("Game 2 should not be solved yet")
	}
	if activeRunAfter2.Games[2].NumGuesses != 2 {
		t.Errorf("Game 2 should have 2 guesses, got %d", activeRunAfter2.Games[2].NumGuesses)
	}

	// Step 4: Submit correct guess for game 2
	guesses[2] = activeRun.Games[2].Answer // Correct for game 2
	// Games 3+ continue using DummyGuess (once DummyGuess is used, it's always used)
	guessesBody, _ = json.Marshal(guessesReq)
	guessesResp, err = client.Post(ts.URL+"/api/guesses", "application/json", bytes.NewBuffer(guessesBody))
	if err != nil {
		t.Fatalf("Failed to call /api/guesses: %v", err)
	}
	defer guessesResp.Body.Close()

	activeRunAfter3, err := storage.GetActiveRun(teamID, runID)
	if err != nil {
		t.Fatalf("Failed to get active run: %v", err)
	}

	// Verify game 2 is now solved with NumGuesses = 3
	if !activeRunAfter3.Games[2].Solved {
		t.Error("Game 2 should be solved after correct guess")
	}
	if activeRunAfter3.Games[2].NumGuesses != 3 {
		t.Errorf("Game 2 should have 3 guesses, got %d", activeRunAfter3.Games[2].NumGuesses)
	}

	// Verify games 3+ are marked as solved when DummyGuess is used
	// Note: DummyGuess indicates the game is already solved (by middleware/client)
	for i := 3; i < len(activeRunAfter3.Games) && i < 10; i++ {
		if !activeRunAfter3.Games[i].Solved {
			t.Errorf("Game %d should be solved (DummyGuess marks it as solved)", i)
		}
		if activeRunAfter3.Games[i].NumGuesses != 0 {
			t.Errorf("Game %d should have 0 guesses (DummyGuess doesn't increment), got %d", i, activeRunAfter3.Games[i].NumGuesses)
		}
	}

	// Step 5: Submit another round - already-solved games should use DummyGuess (as middleware would)
	guesses[0] = common.DummyGuess // Middleware sends DummyGuess for already-solved game 0
	guesses[1] = common.DummyGuess // Middleware sends DummyGuess for already-solved game 1
	// Games 3+ continue using DummyGuess (once DummyGuess is used, it's always used)

	guessesBody, _ = json.Marshal(guessesReq)
	guessesResp, err = client.Post(ts.URL+"/api/guesses", "application/json", bytes.NewBuffer(guessesBody))
	if err != nil {
		t.Fatalf("Failed to call /api/guesses: %v", err)
	}
	defer guessesResp.Body.Close()

	activeRunAfter4, err := storage.GetActiveRun(teamID, runID)
	if err != nil {
		t.Fatalf("Failed to get active run: %v", err)
	}

	// Verify games 0 and 1 are still solved and NumGuesses didn't increment
	if !activeRunAfter4.Games[0].Solved {
		t.Error("Game 0 should remain solved")
	}
	if activeRunAfter4.Games[0].NumGuesses != 2 {
		t.Errorf("Game 0 NumGuesses should remain 2 (already solved), got %d", activeRunAfter4.Games[0].NumGuesses)
	}

	if !activeRunAfter4.Games[1].Solved {
		t.Error("Game 1 should remain solved")
	}
	if activeRunAfter4.Games[1].NumGuesses != 2 {
		t.Errorf("Game 1 NumGuesses should remain 2 (already solved), got %d", activeRunAfter4.Games[1].NumGuesses)
	}

	// Step 6: Call POST /end once at the end to complete the run (all games are now solved)
	endReq := handlers.EndRequest{
		TeamID: teamID,
		RunID:  runID,
	}
	endBody, err := json.Marshal(endReq)
	if err != nil {
		t.Fatalf("Failed to marshal end request: %v", err)
	}

	endResp, err := client.Post(ts.URL+"/end", "application/json", bytes.NewBuffer(endBody))
	if err != nil {
		t.Fatalf("Failed to call /end: %v", err)
	}
	defer endResp.Body.Close()

	if endResp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200 OK, got %d", endResp.StatusCode)
	}

	var endResponse handlers.EndResponse
	if err := json.NewDecoder(endResp.Body).Decode(&endResponse); err != nil {
		t.Fatalf("Failed to decode end response: %v", err)
	}

	// Verify response values
	if !endResponse.Solved {
		t.Error("End response should indicate all games are solved")
	}
	if endResponse.Score <= 0 {
		t.Errorf("End response score should be positive, got %f", endResponse.Score)
	}
	if endResponse.AverageGuesses <= 0 {
		t.Errorf("End response averageGuesses should be positive, got %f", endResponse.AverageGuesses)
	}
	if endResponse.Score != endResponse.AverageGuesses {
		t.Errorf("End response score and averageGuesses should be equal, got score=%f, averageGuesses=%f", endResponse.Score, endResponse.AverageGuesses)
	}

	// Step 7: Verify the correct entry is created in Scores database
	scoreItem, err := storage.GetScore(teamID)
	if err != nil {
		t.Fatalf("Failed to get score: %v", err)
	}

	if scoreItem.TeamID != teamID {
		t.Errorf("ScoreItem TeamID should be %s, got %s", teamID, scoreItem.TeamID)
	}

	if len(scoreItem.CompletedRuns) != 1 {
		t.Fatalf("Expected 1 completed run, got %d", len(scoreItem.CompletedRuns))
	}

	completedRun := scoreItem.CompletedRuns[0]
	if completedRun.RunID != runID {
		t.Errorf("CompletedRun RunID should be %s, got %s", runID, completedRun.RunID)
	}
	if !completedRun.Solved {
		t.Error("CompletedRun should be marked as solved")
	}
	if completedRun.Score != endResponse.Score {
		t.Errorf("CompletedRun Score should be %f, got %f", endResponse.Score, completedRun.Score)
	}
	if completedRun.AverageGuesses != endResponse.AverageGuesses {
		t.Errorf("CompletedRun AverageGuesses should be %f, got %f", endResponse.AverageGuesses, completedRun.AverageGuesses)
	}

	// Step 8: Verify the ActiveRun entry has been removed
	_, err = storage.GetActiveRun(teamID, runID)
	if err == nil {
		t.Fatal("ActiveRun should have been removed after calling /end")
	}
	if err != nil && !strings.Contains(err.Error(), "expired or not found") {
		t.Fatalf("Expected 'expired or not found' error, got: %v", err)
	}
}
