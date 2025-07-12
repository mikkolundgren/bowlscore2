package main

import (
	"log"
	"net/http"
	"os"

	db "github.com/mikkolundgren/bowlscore/database"
	routes "github.com/mikkolundgren/bowlscore/routes"
)

func main() {
	// Initialize database
	dbPath, isPresent := os.LookupEnv("DB_PATH")
	if !isPresent {
		log.Println("No DB_PATH env. Using defaults...")	
		dbPath = "./bowling_scores.db"
	}
	log.Printf("Initializing DB with path: %s", dbPath)
	db.InitDB(dbPath)
	
	defer db.Close()

	// Setup routes
	r := routes.SetupRoutes()

	log.Println("HTTPS server starting on :8080")

	log.Fatal(http.ListenAndServe(":8080", r))
	//log.Fatal(http.ListenAndServeTLS(":8443", "certs/bowlscore.com+1.pem", "certs/bowlscore.com+1-key.pem", r))
}
