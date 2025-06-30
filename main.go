package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type BowlingScore struct {
	ID         int       `json:"id"`
	PlayerID   string    `json:"player_id"`
	Frames     string    `json:"frames"` // JSON string of frame data
	TotalScore int       `json:"total_score"`
	Timestamp  time.Time `json:"timestamp"`
}

type ScoreRequest struct {
	PlayerID   string `json:"player_id"`
	Frames     string `json:"frames"`
	TotalScore int    `json:"total_score"`
}

var db *sql.DB

func main() {
	// Initialize database
	var err error
	db, err = sql.Open("sqlite3", "./bowling_scores.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create table if not exists
	createTable()

	// Setup routes
	r := mux.NewRouter()

	// API routes
	r.HandleFunc("/api/scores", saveScore).Methods("POST")
	r.HandleFunc("/api/scores", listScores).Methods("GET")
	r.HandleFunc("/api/scores/{id}", deleteScore).Methods("DELETE")
	r.HandleFunc("/api/player-progress", getPlayerProgress).Methods("GET")

	// Serve static files
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	// Enable CORS
	r.Use(corsMiddleware)

	log.Println("HTTPS server starting on :8443")
	log.Println("Access the application at: https://localhost:8443")
	log.Fatal(http.ListenAndServeTLS(":8443", "certs/bowlscore.com+1.pem", "certs/bowlscore.com+1-key.pem", r))
}

func createTable() {
	query := `
	CREATE TABLE IF NOT EXISTS bowling_scores (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		player_id TEXT NOT NULL,
		frames TEXT NOT NULL,
		total_score INTEGER NOT NULL,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func saveScore(w http.ResponseWriter, r *http.Request) {
	var scoreReq ScoreRequest
	if err := json.NewDecoder(r.Body).Decode(&scoreReq); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	query := `
	INSERT INTO bowling_scores (player_id, frames, total_score, timestamp)
	VALUES (?, ?, ?, ?)`

	result, err := db.Exec(query, scoreReq.PlayerID, scoreReq.Frames, scoreReq.TotalScore, time.Now())
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

func listScores(w http.ResponseWriter, r *http.Request) {
	query := `
	SELECT id, player_id, frames, total_score, timestamp
	FROM bowling_scores
	ORDER BY timestamp DESC`

	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Error querying scores: %v", err)
		http.Error(w, "Failed to retrieve scores", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var scores []BowlingScore
	for rows.Next() {
		var score BowlingScore
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

func deleteScore(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM bowling_scores WHERE id = ?`
	result, err := db.Exec(query, id)
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

func getPlayerProgress(w http.ResponseWriter, r *http.Request) {
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

	rows, err := db.Query(query, playerID)
	if err != nil {
		log.Printf("Error querying player progress: %v", err)
		http.Error(w, "Failed to retrieve player progress", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type DailyProgress struct {
		Date        string  `json:"date"`
		Average     float64 `json:"average"`
		GamesPlayed int     `json:"games_played"`
	}

	var progress []DailyProgress
	for rows.Next() {
		var p DailyProgress
		err := rows.Scan(&p.Date, &p.Average, &p.GamesPlayed)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
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
