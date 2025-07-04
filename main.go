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
		dbPath = "./bowling_scores.db"
	}
	db.InitDB(dbPath)
	defer db.Close()

	// Setup routes
	r := routes.SetupRoutes()

	log.Println("HTTPS server starting on :8080")

	log.Fatal(http.ListenAndServe(":8080", r))
	//log.Fatal(http.ListenAndServeTLS(":8443", "certs/bowlscore.com+1.pem", "certs/bowlscore.com+1-key.pem", r))
}
