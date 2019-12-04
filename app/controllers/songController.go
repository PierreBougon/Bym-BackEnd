package controllers

import (
	"github.com/PierreBougon/Bym-BackEnd/app/models"
	u "github.com/PierreBougon/Bym-BackEnd/app/utils"
	"github.com/PierreBougon/Bym-BackEnd/app/websocket"
	"github.com/gorilla/mux"

	"encoding/json"
	"net/http"
	"strconv"
)

var CreateSong = func(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(uint)
	var song = &models.Song{}
	err := json.NewDecoder(r.Body).Decode(song) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := song.Create(user, websocket.NotifyPlaylistSubscribers, websocket.PlaylistNeedRefresh)
	if resp["status"] == false {
		w.WriteHeader(http.StatusBadRequest)
	}
	u.Respond(w, resp)
}

var GetSongs = func(w http.ResponseWriter, r *http.Request) {
	var param string

	user := r.Context().Value("user").(uint)
	vals := r.URL.Query()                 // Returns a url.Values, which is a map[string][]string
	playlistId, ok := vals["playlist_id"] // Note type, not ID. ID wasn't specified anywhere.

	if ok && len(playlistId) >= 1 {
		param = playlistId[0] // The first `?type=model`
	} else {
		w.WriteHeader(http.StatusBadRequest)
		u.Respond(w, u.Message(false, "Invalid request need playlist id"))
		return
	}

	plistid, err := strconv.Atoi(param)
	if err != nil {
		//The passed path parameter is not an integer
		w.WriteHeader(http.StatusBadRequest)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	data := models.GetSongs(uint(plistid), user)
	if data == nil {
		u.RespondBadRequestWithMessage(w, "Invalid request, playlist Id doesn't match with any playlist")
		return
	}
	resp := u.Message(true, "success")
	resp["songs"] = data
	u.Respond(w, resp)
}

var UpdateSong = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 32)
	if err != nil {
		u.RespondBadRequest(w)
		return
	}
	user := r.Context().Value("user").(uint)
	var song = &models.Song{}
	err = json.NewDecoder(r.Body).Decode(song) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.RespondBadRequest(w)
		return
	}
	resp := song.UpdateSong(user, uint(id), song, websocket.NotifyPlaylistSubscribers, websocket.PlaylistNeedRefresh)
	if resp["status"] == false {
		w.WriteHeader(http.StatusBadRequest)
	}
	u.Respond(w, resp)
}

var DeleteSong = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 32)
	if err != nil {
		u.RespondBadRequest(w)
		return
	}
	user := r.Context().Value("user").(uint)
	resp := (&models.Song{}).DeleteSong(user, uint(id), websocket.NotifyPlaylistSubscribers, websocket.PlaylistNeedRefresh)
	if resp["status"] == false {
		w.WriteHeader(http.StatusBadRequest)
	}
	u.Respond(w, resp)
}
