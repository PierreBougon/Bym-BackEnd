package controllers

import (
	"github.com/PierreBougon/Bym-BackEnd/models"
	u "github.com/PierreBougon/Bym-BackEnd/utils"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// Not directly related to a model it's an interface for rankings conatained in the Song model

var GetRankings = func(w http.ResponseWriter, r *http.Request) {
	var param string

	vals := r.URL.Query()                  // Returns a url.Values, which is a map[string][]string
	playlist_id, ok := vals["playlist_id"] // Note type, not ID. ID wasn't specified anywhere.

	if ok && len(playlist_id) >= 1 {
		param = playlist_id[0] // The first `?type=model`
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

	data := models.GetSongsRanking(uint(plistid))
	resp := u.Message(true, "success")
	resp["songs_rankings"] = data
	u.Respond(w, resp)
}

var GetRanking = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 32)
	if err != nil {
		u.RespondBadRequest(w)
		return
	}

	data := models.GetSongRankingById(uint(id))
	resp := u.Message(true, "success")
	resp["song_ranking"] = data
	u.Respond(w, resp)
}
