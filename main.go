package main

import (
	"github.com/PierreBougon/Bym-BackEnd/app"
	"github.com/PierreBougon/Bym-BackEnd/controllers"
	"github.com/PierreBougon/Bym-BackEnd/moesif"

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
		port = "443"
	}

	fmt.Println(port)

	// API
	api := router.PathPrefix("/api").Subrouter()
	api.Use(moesif.MiddlewareWrapper)

	//		Auth
	auth := api.PathPrefix("/user").Subrouter()
	auth.HandleFunc("/new", controllers.CreateAccount).Methods("POST")
	auth.HandleFunc("/login", controllers.Authenticate).Methods("POST")
	auth.HandleFunc("/update_password", controllers.UpdatePassword).Methods("PATCH")

	//		Playlist
	playlist := api.PathPrefix("/playlist").Subrouter()
	playlist.HandleFunc("", controllers.CreatePlaylist).Methods("POST")
	playlist.HandleFunc("", controllers.GetPlaylists).Methods("GET")
	playlist.HandleFunc("/{id}", controllers.GetPlaylist).Methods("GET")
	playlist.HandleFunc("/{id}", controllers.UpdatePlaylist).Methods("PUT")
	playlist.HandleFunc("/{id}", controllers.DeletePlaylist).Methods("DELETE")

	//		Songs
	song := api.PathPrefix("/song").Subrouter()
	song.HandleFunc("", controllers.CreateSong).Methods("POST")
	song.HandleFunc("", controllers.GetSongs).Methods("GET")
	song.HandleFunc("/{id}", controllers.UpdateSong).Methods("PUT")
	song.HandleFunc("/{id}", controllers.DeleteSong).Methods("DELETE")

	//			Ranking (Fraction of Songs parsed to access it directly)
	ranking := song.PathPrefix("/ranking").Subrouter()
	ranking.HandleFunc("", controllers.GetRankings).Methods("GET")
	ranking.HandleFunc("/{id}", controllers.GetRanking).Methods("GET")

	//		Votes
	vote := api.PathPrefix("/vote").Subrouter()
	vote.HandleFunc("", controllers.UpdateOrCreateVote).Methods("PUT")
	vote.HandleFunc("", controllers.GetVote).Methods("GET")
	//	vote.HandleFunc("/{id}", controllers.DeleteSong).Methods("DELETE")

	err := http.ListenAndServe(":"+port, router) //Launch the app, visit localhost:8000/api
	if err != nil {
		fmt.Print(err)
	}
}
