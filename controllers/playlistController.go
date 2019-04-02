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
		u.RespondBadRequest(w)
		return
	}

	resp := playlist.Create(user)
	u.Respond(w, resp)
}

var GetPlaylists = func(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(uint)

	data := models.GetPlaylists(uint(user))
	resp := u.Message(true, "success")
	resp["playlists"] = data
	u.Respond(w, resp)
}

var GetPlaylist = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		u.RespondBadRequest(w)
		return
	}

	data := models.GetPlaylist(uint(id))
	resp := u.Message(true, "success")
	resp["playlist"] = data
	u.Respond(w, resp)
}

var UpdatePlaylist = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 32)
	if err != nil {
		u.RespondBadRequest(w)
		return
	}
	user := r.Context().Value("user").(uint)
	var playlist = &models.Playlist{}
	err = json.NewDecoder(r.Body).Decode(playlist) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.RespondBadRequest(w)
		return
	}
	resp := playlist.UpdatePlaylist(user, uint(id), playlist)
	if resp["status"] == false {
		w.WriteHeader(http.StatusBadRequest)
	}
	u.Respond(w, resp)
}

var DeletePlaylist = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 32)
	if err != nil {
		u.RespondBadRequest(w)
		return
	}
	user := r.Context().Value("user").(uint)
	resp := (&models.Playlist{}).DeletePlaylist(user, uint(id))
	if resp["status"] == false {
		w.WriteHeader(http.StatusBadRequest)
	}
	u.Respond(w, resp)
}
