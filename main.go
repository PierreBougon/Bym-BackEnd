package main

import (
	"github.com/PierreBougon/Bym-BackEnd/app"
	"github.com/PierreBougon/Bym-BackEnd/controllers"

	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()
	router.Use(app.JwtAuthentication) //attach JWT auth middleware

	port := os.Getenv("PORT") //Get port from .env file, we did not specify any port so this should return an empty string when tested locally
	if port == "" {
		port = "443" //localhost
	}

	fmt.Println(port)

	// Auth
	router.HandleFunc("/api/user/new", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login", controllers.Authenticate).Methods("POST")
	router.HandleFunc("/api/user/update_password", controllers.UpdatePassword).Methods("PATCH")

	// Playlist
	router.HandleFunc("/api/playlist", controllers.CreatePlaylist).Methods("POST")
	router.HandleFunc("/api/playlist", controllers.GetPlaylists).Methods("GET")
	router.HandleFunc("/api/playlist/{id}", controllers.GetPlaylist).Methods("GET")
	router.HandleFunc("/api/playlist/{id}", controllers.UpdatePlaylist).Methods("PUT")
	router.HandleFunc("/api/playlist/{id}", controllers.DeletePlaylist).Methods("DELETE")

	// Songs
	router.HandleFunc("/api/song", controllers.CreateSong).Methods("POST")
	router.HandleFunc("/api/song", controllers.GetSongs).Methods("GET")
	router.HandleFunc("/api/song/{id}", controllers.UpdateSong).Methods("PUT")
	router.HandleFunc("/api/song/{id}", controllers.DeleteSong).Methods("DELETE")

	// Ranking (Fraction of Songs parsed to access it directly)
	router.HandleFunc("/api/song/ranking", controllers.GetRankings).Methods("GET")
	router.HandleFunc("/api/song/ranking/{id}", controllers.GetRanking).Methods("GET")

	// Votes
	router.HandleFunc("/api/vote", controllers.UpdateOrCreateVote).Methods("PUT")
	router.HandleFunc("/api/vote", controllers.GetVote).Methods("GET")
	//	router.HandleFunc("/api/vote/{id}", controllers.DeleteSong).Methods("DELETE")

	router.Use(app.JwtAuthentication)

	err := http.ListenAndServe(":"+port, router) //Launch the app, visit localhost:8000/api
	if err != nil {
		fmt.Print(err)
	}
}
