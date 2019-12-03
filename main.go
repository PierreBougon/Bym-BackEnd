package main

import (
	"github.com/PierreBougon/Bym-BackEnd/app/auth"
	"github.com/PierreBougon/Bym-BackEnd/app/controllers"
	"github.com/PierreBougon/Bym-BackEnd/app/moesif"
	u "github.com/PierreBougon/Bym-BackEnd/app/utils"

	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {

	// Get port from .env file, we did not specify any port so this should return an empty string when tested locally
	port := os.Getenv("PORT")
	if port == "" {
		port = "443"
	}

	router := newAPIRouter(port)

	//Launch the app, visit localhost:443/api
	err := http.ListenAndServe(":"+port, router)

	if err != nil {
		fmt.Print(err)
	}
}

func newAPIRouter(port string) *mux.Router {
	router := mux.NewRouter()
	router.Use(auth.JwtAuthentication) //attach JWT auth middleware

	fmt.Println(port)

	// API
	api := router.PathPrefix("/api").Subrouter()
	api.Use(moesif.MiddlewareWrapper)

	// Respond a basic success if anyone wants to get from / or /api to let them now the url is correct and server is up
	router.HandleFunc("", u.RespondBasicSuccess).Methods("GET")
	router.HandleFunc("/", u.RespondBasicSuccess).Methods("GET")
	api.HandleFunc("", u.RespondBasicSuccess).Methods("GET")

	// 		Connect Websocket
	router.HandleFunc("/ws", controllers.ConnectWebSocket).Methods("GET")
	//		Auth / Account
	attachAuthRoutes(api)
	//		Playlist
	attachPlaylistRoutes(api)
	//		Songs
	attachSongRoutes(api)
	//		Ranking (Fraction of Songs parsed to access it directly)
	attachRankingRoutes(api)
	//		Votes
	attachVoteRoutes(api)

	return router
}

func attachAuthRoutes(api *mux.Router) {
	auth := api.PathPrefix("/user").Subrouter()
	auth.HandleFunc("/new", controllers.CreateAccount).Methods("POST")
	auth.HandleFunc("/login", controllers.Authenticate).Methods("POST")
	auth.HandleFunc("/delete", controllers.DeleteAccount).Methods("DELETE")
	auth.HandleFunc("", controllers.UpdateAccount).Methods("PUT")
	auth.HandleFunc("", controllers.GetAccount).Methods("GET")
	auth.HandleFunc("/update_password", controllers.UpdatePassword).Methods("PATCH")
}

func attachPlaylistRoutes(api *mux.Router) {
	playlist := api.PathPrefix("/playlist").Subrouter()
	playlist.HandleFunc("", controllers.CreatePlaylist).Methods("POST")
	playlist.HandleFunc("", controllers.GetPlaylists).Methods("GET")
	playlist.HandleFunc("/{id}", controllers.GetPlaylist).Methods("GET")
	playlist.HandleFunc("/{id}", controllers.UpdatePlaylist).Methods("PUT")
	playlist.HandleFunc("/{id}", controllers.DeletePlaylist).Methods("DELETE")
	playlist.HandleFunc("/join/{id}", controllers.JoinPlaylist).Methods("POST")
	playlist.HandleFunc("/leave/{id}", controllers.LeavePlaylist).Methods("DELETE")
	playlist.HandleFunc("/change_user_acl/{id}", controllers.ChangeAclOnPlaylist).Methods("POST")
	playlist.HandleFunc("/get_role/{id}", controllers.GetPlaylistRole).Methods("GET")
}

func attachSongRoutes(api *mux.Router) {
	song := api.PathPrefix("/song").Subrouter()
	song.HandleFunc("", controllers.CreateSong).Methods("POST")
	song.HandleFunc("", controllers.GetSongs).Methods("GET")
	song.HandleFunc("/{id}", controllers.UpdateSong).Methods("PUT")
	song.HandleFunc("/{id}", controllers.DeleteSong).Methods("DELETE")
}

func attachRankingRoutes(api *mux.Router) {
	ranking := api.PathPrefix("/song/ranking").Subrouter()
	ranking.HandleFunc("", controllers.GetRankings).Methods("GET")
	ranking.HandleFunc("/{id}", controllers.GetRanking).Methods("GET")
}

func attachVoteRoutes(api *mux.Router) {
	vote := api.PathPrefix("/vote").Subrouter()
	vote.HandleFunc("", controllers.UpdateOrCreateVote).Methods("PUT")
	vote.HandleFunc("", controllers.GetVote).Methods("GET")
	//	vote.HandleFunc("/{id}", controllers.DeleteSong).Methods("DELETE")
}
