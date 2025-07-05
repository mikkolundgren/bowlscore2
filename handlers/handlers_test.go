package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gorilla/mux"
	_ "modernc.org/sqlite"
	"github.com/stretchr/testify/assert"
	
	models "github.com/mikkolundgren/bowlscore/models"
	database "github.com/mikkolundgren/bowlscore/database"
)

const testDBPath = "./test_bowling_scores.db"

func TestMain(m *testing.M) {
	setup()
	defer database.Close()
	m.Run()
	
	teardownTestDB()
	
}

func setup () {
	// Remove any existing test database
	os.Remove(testDBPath)
	database.InitDB(testDBPath)
}

func teardownTestDB() {
	err := os.Remove(testDBPath)
	if err != nil {
		fmt.Println("Error removing test database.")
	}
}

func TestSaveScore(t *testing.T) {

	scoreReq := models.ScoreRequest{
		PlayerID:   "player1",
		Frames:     `[10,7,3,9,0,10,0,8,8,2,0,6,10,10,10,8,1]`,
		TotalScore: 167,
	}
	body, _ := json.Marshal(scoreReq)
	req, _ := http.NewRequest("POST", "/api/scores", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	SaveScore(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]any
	json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NotNil(t, response["id"])
	assert.Equal(t, "Score saved successfully", response["message"])
}

func TestListScores(t *testing.T) {

	// Insert test data
	_, err := database.Db.Exec(`INSERT INTO bowling_scores (player_id, frames, total_score) VALUES (?, ?, ?)`,
		"player_insert", `[10,7,3,9,0,10,0,8,8,2,0,6,10,10,10,8,1]`, 167)
	if err != nil {
		t.Fatal("Failed to insert test data:", err)
	}

	req, _ := http.NewRequest("GET", "/api/scores", nil)
	rr := httptest.NewRecorder()

	ListScores(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var scores []models.BowlingScore
	json.Unmarshal(rr.Body.Bytes(), &scores)
	assert.Equal(t, 2, len(scores))
	assert.Equal(t, "player_insert", scores[1].PlayerID)
	assert.Equal(t, 167, scores[1].TotalScore)
}

func TestDeleteScore(t *testing.T) {
	
	// Insert test data
	result, err := database.Db.Exec(`INSERT INTO bowling_scores (player_id, frames, total_score) VALUES (?, ?, ?)`,
		"player6", `[10,7,3,9,0,10,0,8,8,2,0,6,10,10,10,8,1]`, 167)
	if err != nil {
		t.Fatal("Failed to insert test data:", err)
	}
	id, _ := result.LastInsertId()

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/scores/%d", id), nil)
	rr := httptest.NewRecorder()

	// Set up path variables
	req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprintf("%d", id)})

	DeleteScore(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]string
	json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Equal(t, "Score deleted successfully", response["message"])
}

func TestGetPlayerProgress(t *testing.T) {

	// Insert test data with different dates
	testData := []struct {
		playerID  string
		frames    string
		score     int
		timestamp time.Time
	}{
		{"player2", `[10,7,3,9,0,10,0,8,8,2,0,6,10,10,10,8,1]`, 167, time.Now().AddDate(0, 0, -2)},
		{"player2", `[9,0,8,1,0,10,10,10,8,1,9,0,7,3,10,8,1]`, 145, time.Now().AddDate(0, 0, -2)},
		{"player2", `[10,10,10,10,10,10,10,10,10,10,10,10]`, 300, time.Now().AddDate(0, 0, -1)},
	}

	for _, data := range testData {
		_, err := database.Db.Exec(`INSERT INTO bowling_scores (player_id, frames, total_score, timestamp) VALUES (?, ?, ?, ?)`,
			data.playerID, data.frames, data.score, data.timestamp.Format("2006-01-02 15:04:05"))
		if err != nil {
			t.Fatal("Failed to insert test data:", err)
		}
	}

	req, _ := http.NewRequest("GET", "/api/player-progress?player_id=player2", nil)
	rr := httptest.NewRecorder()

	GetPlayerProgress(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var progress []models.DailyProgress
	json.Unmarshal(rr.Body.Bytes(), &progress)
	assert.Equal(t, 2, len(progress)) // Should have 2 days of data

	// Find today-2 data
	var day1 models.DailyProgress
	for _, p := range progress {
		if p.GamesPlayed == 2 {
			day1 = p
		}
	}

	assert.Equal(t, 2, day1.GamesPlayed)
	assert.Equal(t, 156.0, day1.Average) // (167 + 145) / 2 = 156
}
