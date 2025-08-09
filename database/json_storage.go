package database

import (
	"encoding/json"
	"os"
	"sort"
	"sync"
	"time"

	models "github.com/mikkolundgren/bowlscore/models"
)

var JSONDb *JSONStorage

func (js *JSONStorage) DeleteScore(id int) error {
	scores, err := js.readScores()
	if err != nil {
		return err
	}

	found := false
	newScores := make([]models.BowlingScore, 0, len(scores))
	for _, score := range scores {
		if score.ID != id {
			newScores = append(newScores, score)
		} else {
			found = true
		}
	}

	if !found {
		return os.ErrNotExist
	}

	return js.writeScores(newScores)
}

type JSONStorage struct {
	filePath string
	mu       sync.Mutex
}

func Init(filePath string) {
	// Create empty file if it doesn't exist
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		os.WriteFile(filePath, []byte("[]"), 0644)
	}
}

func NewJSONStorage() *JSONStorage {
	dbPath, _ := os.LookupEnv("DB_PATH")
	return &JSONStorage{
		filePath: dbPath,
	}
}

func (js *JSONStorage) readScores() ([]models.BowlingScore, error) {
	js.mu.Lock()
	defer js.mu.Unlock()

	data, err := os.ReadFile(js.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []models.BowlingScore{}, nil
		}
		return nil, err
	}

	var scores []models.BowlingScore
	err = json.Unmarshal(data, &scores)
	return scores, err
}

func (js *JSONStorage) writeScores(scores []models.BowlingScore) error {
	js.mu.Lock()
	defer js.mu.Unlock()

	data, err := json.MarshalIndent(scores, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(js.filePath, data, 0644)
}

func (js *JSONStorage) SaveScore(playerID, frames string, totalScore int) (int, error) {
	scores, err := js.readScores()
	if err != nil {
		return 0, err
	}

	newID := 1
	if len(scores) > 0 {
		newID = scores[len(scores)-1].ID + 1
	}
	now := time.Now()
	dateTime := now.Format("2006-01-02 15:04:05")

	newScore := models.BowlingScore{
		ID:         newID,
		PlayerID:   playerID,
		Frames:     frames,
		TotalScore: totalScore,
		Timestamp:  dateTime,
	}

	scores = append(scores, newScore)
	err = js.writeScores(scores)
	return newID, err
}

func (js *JSONStorage) ListScores() ([]models.BowlingScore, error) {
	return js.readScores()
}

func (js *JSONStorage) GetPlayerProgress(playerID string) (map[string]struct {
	Average     float64
	GamesPlayed int
}, error) {
	scores, err := js.readScores()
	if err != nil {
		return nil, err
	}

	progress := make(map[string]struct {
		Average     float64
		GamesPlayed int
	})

	for _, score := range scores {
		if score.PlayerID == playerID {
			date := score.Timestamp[:10] // Extract YYYY-MM-DD
			entry := progress[date]
			entry.GamesPlayed++
			entry.Average = (entry.Average*float64(entry.GamesPlayed-1) + float64(score.TotalScore)) / float64(entry.GamesPlayed)
			progress[date] = entry
		}
	}

	// Sort dates chronologically
	keys := make([]string, 0, len(progress))
	for k := range progress {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	sortedProgress := make(map[string]struct {
		Average     float64
		GamesPlayed int
	}, len(progress))
	for _, k := range keys {
		sortedProgress[k] = progress[k]
	}

	return sortedProgress, nil
}
