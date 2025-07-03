package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	handlers "github.com/mikkolundgren/bowlscore/handlers"
	middleware "github.com/mikkolundgren/bowlscore/middleware"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	// API routes
	r.HandleFunc("/api/scores", handlers.SaveScore).Methods("POST")
	r.HandleFunc("/api/scores", handlers.ListScores).Methods("GET")
	r.HandleFunc("/api/scores/{id}", handlers.DeleteScore).Methods("DELETE")
	r.HandleFunc("/api/player-progress", handlers.GetPlayerProgress).Methods("GET")

	// Serve static files
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	// Enable CORS
	r.Use(middleware.CorsMiddleware)

	return r
}
