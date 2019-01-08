package controllers

import (
	"Bym-BackEnd/models"
	u "Bym-BackEnd/utils"
	"encoding/json"
	"net/http"
	"strconv"
)

var CreateSong = func(w http.ResponseWriter, r *http.Request) {
	var song = &models.Song{}
	err := json.NewDecoder(r.Body).Decode(song) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := song.Create()
	u.Respond(w, resp)
}

var GetSongs = func(w http.ResponseWriter, r *http.Request) {
	var param string

	vals := r.URL.Query()                  // Returns a url.Values, which is a map[string][]string
	playlist_id, ok := vals["playlist_id"] // Note type, not ID. ID wasn't specified anywhere.

	if ok && len(playlist_id) >= 1 {
		param = playlist_id[0] // The first `?type=model`
	} else {
		u.Respond(w, u.Message(false, "Invalid request need playlist id"))
		return
	}

	plistid, err := strconv.Atoi(param)
	if err != nil {
		//The passed path parameter is not an integer
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	data := models.GetSongs(uint(plistid))
	resp := u.Message(true, "success")
	resp["songs"] = data
	u.Respond(w, resp)
}
