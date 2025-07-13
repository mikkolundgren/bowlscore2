package database

import (
	"database/sql"
	"log"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

var Db *sql.DB
var dbPath string

func InitDB(dbPathParam string) {
	dbPath = dbPathParam
	createTable(dbPath)
	
}

func Open() {
	
	var err error
	Db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		log.Println("error opening db ", dbPath, err)
	}	
}

func Close() {
	Db.Close()
}

func createTable(dbPath string)  {
	
	var err error
	
	log.Printf("creating table...")
	query := `
	CREATE TABLE IF NOT EXISTS bowling_scores (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		player_id TEXT NOT NULL,
		frames TEXT NOT NULL,
		total_score INTEGER NOT NULL,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	Db, err = sql.Open("sqlite", dbPath)
	i := 0
	for i < 5 {
		_, err = Db.Exec(query)
		if err != nil {
			log.Println("error creating table", err)
			if strings.Contains(err.Error(), "SQLITE_BUSY") {
				log.Println("retrying table creation...")
				time.Sleep(3 * time.Second)
			}
			i = i + 1 
		} else {
			break
		}
	}
	defer Db.Close()
}
