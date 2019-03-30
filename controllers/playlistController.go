package controllers

import (
	"github.com/PierreBougon/Bym-BackEnd/models"
	u "github.com/PierreBougon/Bym-BackEnd/utils"

	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

var CreatePlaylist = func(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value("user").(uint)
	var playlist = &models.Playlist{}
	err := json.NewDecoder(r.Body).Decode(playlist) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := playlist.Create(user)
	u.Respond(w, resp)
}

var GetPlaylists = func(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(uint)

	data := models.GetPlaylistByUser(uint(user))
	resp := u.Message(true, "success")
	resp["playlists"] = data
	u.Respond(w, resp)
}

var GetPlaylist = func(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		//The passed path parameter is not an integer
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	data := models.GetPlaylistById(uint(id))
	resp := u.Message(true, "success")
	resp["playlist"] = data
	u.Respond(w, resp)
}
