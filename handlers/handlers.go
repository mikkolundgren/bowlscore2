package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	models "github.com/mikkolundgren/bowlscore/models"
	db "github.com/mikkolundgren/bowlscore/database"
)

func SaveScore(w http.ResponseWriter, r *http.Request) {
	var scoreReq models.ScoreRequest
	if err := json.NewDecoder(r.Body).Decode(&scoreReq); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	query := `
	INSERT INTO bowling_scores (player_id, frames, total_score, timestamp)
	VALUES (?, ?, ?, ?)`

	now := time.Now()
	dateTime := now.Format("2006-01-02 15:04:05")
	
	result, err := db.Db.Exec(query, scoreReq.PlayerID, scoreReq.Frames, scoreReq.TotalScore, dateTime)
	if err != nil {
		log.Printf("Error saving score: %v", err)
		http.Error(w, "Failed to save score", http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()

	response := map[string]any{
		"id":      id,
		"message": "Score saved successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func ListScores(w http.ResponseWriter, r *http.Request) {
	query := `
	SELECT id, player_id, frames, total_score, timestamp
	FROM bowling_scores
	ORDER BY timestamp DESC`

	rows, err := db.Db.Query(query)
	if err != nil {
		log.Printf("Error querying scores: %v", err)
		http.Error(w, "Failed to retrieve scores", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var scores []models.BowlingScore
	for rows.Next() {
		var score models.BowlingScore
		err := rows.Scan(&score.ID, &score.PlayerID, &score.Frames, &score.TotalScore, &score.Timestamp)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		scores = append(scores, score)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(scores)
}

func DeleteScore(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM bowling_scores WHERE id = ?`
	result, err := db.Db.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting score: %v", err)
		http.Error(w, "Failed to delete score", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Score not found", http.StatusNotFound)
		return
	}

	response := map[string]string{
		"message": "Score deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetPlayerProgress(w http.ResponseWriter, r *http.Request) {
	playerID := r.URL.Query().Get("player_id")
	if playerID == "" {
		http.Error(w, "Missing player_id parameter", http.StatusBadRequest)
		return
	}

	query := `
    SELECT DATE(timestamp) as date,
           AVG(total_score) as average_score,
           COUNT(*) as games_played
    FROM bowling_scores
    WHERE player_id = ?
    GROUP BY DATE(timestamp)
    ORDER BY date`

	rows, err := db.Db.Query(query, playerID)
	if err != nil {
		log.Printf("Error querying player progress: %v", err)
		http.Error(w, "Failed to retrieve player progress", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var progress []models.DailyProgress
	for rows.Next() {
		var p models.DailyProgress
		err := rows.Scan(&p.Date, &p.Average, &p.GamesPlayed)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		
		// todo: function that calculates the ratios
		progress = append(progress, p)
	}

	if len(progress) == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "No data for player"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(progress)
}
