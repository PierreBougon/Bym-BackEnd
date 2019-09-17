package controllers

import (
	"Bym-BackEnd/models"
	u "Bym-BackEnd/utils"

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
	vals := r.URL.Query()          // Returns a url.Values, which is a map[string][]string
	user_id, ok := vals["user_id"] // Note type, not ID. ID wasn't specified anywhere.

	user := uint(0)
	if ok && len(user_id) >= 1 {
		id, err := strconv.ParseUint(user_id[0], 10, 32) // The first `?type=model`
		if err != nil {
			u.RespondBadRequest(w)
			return
		}
		user = uint(id)
	} else {
		user = r.Context().Value("user").(uint)
	}

	data := models.GetPlaylistsByUser(uint(user))
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

	data := models.GetPlaylistById(uint(id))
	if data == nil {
		u.RespondBadRequest(w)
		return
	}
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
