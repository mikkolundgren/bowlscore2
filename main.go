package main

import (
	"log"
	"net/http"
	db "github.com/mikkolundgren/bowlscore/database"
	routes "github.com/mikkolundgren/bowlscore/routes"
)

func main() {
	// Initialize database
	db.InitDB()
	defer db.Close()

	// Setup routes
	r := routes.SetupRoutes()

	log.Println("HTTPS server starting on :8443")
	log.Println("Access the application at: https://localhost:8443")
	log.Fatal(http.ListenAndServeTLS(":8443", "certs/bowlscore.com+1.pem", "certs/bowlscore.com+1-key.pem", r))
}
