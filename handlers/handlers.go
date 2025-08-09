package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	db "github.com/mikkolundgren/bowlscore/database"
	models "github.com/mikkolundgren/bowlscore/models"
)

func SaveScore(w http.ResponseWriter, r *http.Request) {

	var scoreReq models.ScoreRequest
	if err := json.NewDecoder(r.Body).Decode(&scoreReq); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	js := db.NewJSONStorage()
	// result, err := db.Db.Exec(query, scoreReq.PlayerID, scoreReq.Frames, scoreReq.TotalScore, dateTime)
	id, err := js.SaveScore(scoreReq.PlayerID, scoreReq.Frames, scoreReq.TotalScore)
	if err != nil {
		log.Printf("Error saving score: %v", err)
		http.Error(w, "Failed to save score", http.StatusInternalServerError)
		return
	}

	response := map[string]any{
		"id":      id,
		"message": "Score saved successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func ListScores(w http.ResponseWriter, r *http.Request) {
	js := db.NewJSONStorage()
	rows, err := js.ListScores()
	if err != nil {
		log.Printf("Error querying scores: %v", err)
		http.Error(w, "Failed to retrieve scores", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rows)
}

func DeleteScore(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	js := db.NewJSONStorage()
	err = js.DeleteScore(id)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "Score not found", http.StatusNotFound)
			return
		}
		log.Printf("Error deleting score: %v", err)
		http.Error(w, "Failed to delete score", http.StatusInternalServerError)
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

	js := db.NewJSONStorage()
	scores, err := js.GetPlayerProgress(playerID)
	if err != nil {
		log.Printf("Error getting player progress: %v", err)
		http.Error(w, "Failed to retrieve player progress", http.StatusInternalServerError)
		return
	}

	var progress []models.DailyProgress
	for date, stats := range scores {
		progress = append(progress, models.DailyProgress{
			Date:        date,
			Average:     stats.Average,
			GamesPlayed: stats.GamesPlayed,
		})
	}

	if len(progress) == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "No data for player"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(progress)
}
