package database

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var Db *sql.DB
var dbPath string

func InitDB(dbPathParam string) {
	var err error
	defer Close()
	dbPath = dbPathParam
	Db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}	
	createTable()
}

func Open() {
	
	var err error
	Db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}	
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
