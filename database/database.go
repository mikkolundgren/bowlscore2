package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

func InitDB() {
	var err error
	dbPath, isPresent := os.LookupEnv("DB_PATH")
	if !isPresent {
		dbPath = "./bowling_scores.db"
	}
	Db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	createTable()
}

func Close() {
	Db.Close()
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

	_, err := Db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}
