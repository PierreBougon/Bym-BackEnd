package main

import (
	"Bym-BackEnd/app"
	"Bym-BackEnd/controllers"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

func main() {

	router := mux.NewRouter()
	router.Use(app.JwtAuthentication) //attach JWT auth middleware

	port := os.Getenv("PORT") //Get port from .env file, we did not specify any port so this should return an empty string when tested locally
	if port == "" {
		port = "80" //localhost
	}

	fmt.Println(port)

	// Auth
	router.HandleFunc("/api/user/new", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login", controllers.Authenticate).Methods("POST")

	// Playlist
	router.HandleFunc("/api/playlist", controllers.CreatePlaylist).Methods("POST")
	router.HandleFunc("/api/playlist", controllers.GetPlaylists).Methods("GET")
	router.HandleFunc("/api/playlist/{id}", controllers.GetPlaylist).Methods("GET")

	// Songs
	router.HandleFunc("/api/song", controllers.CreateSong).Methods("POST")
	router.HandleFunc("/api/song", controllers.GetSongs).Methods("GET")

	router.Use(app.JwtAuthentication)

	err := http.ListenAndServe(":"+port, router) //Launch the app, visit localhost:8000/api
	if err != nil {
		fmt.Print(err)
	}
}
