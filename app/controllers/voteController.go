package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/PierreBougon/Bym-BackEnd/app/models"
	u "github.com/PierreBougon/Bym-BackEnd/app/utils"
	"github.com/PierreBougon/Bym-BackEnd/app/websocket"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"net/url"
	"strconv"
)

// Get author vote for a given song
var GetVote = func(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(uint)
	songId, err := getSongId(r.URL.Query())
	if err != nil {
		u.RespondBadRequest(w)
		return
	}

	data := models.GetPersonalVoteBySongId(uint(songId), user)
	resp := u.Message(true, "success")
	resp["vote"] = data
	u.Respond(w, resp)
}

var GetVotes = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	songId, err := strconv.ParseUint(params["song_id"], 10, 32)
	if err != nil {
		u.RespondBadRequest(w)
		return
	}

	data := models.GetVotesBySongId(uint(songId))
	resp := u.Message(true, "success")
	resp["votes"] = data
	u.Respond(w, resp)
}

var UpdateOrCreateVote = func(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(uint)

	songId, err := getSongId(r.URL.Query())
	if err != nil {
		fmt.Println(err)
		u.RespondBadRequest(w)
		return
	}
	vote := &models.Vote{}
	err = json.NewDecoder(r.Body).Decode(&vote) //decode the request body into struct and failed if any error occur
	if err != nil {
		fmt.Println(err)
		u.RespondBadRequest(w)
		return
	}

	resp := u.Message(true, "This request has performed no action")
	if vote.UpVote == true {
		resp = models.UpVoteSong(uint(songId), user, websocket.NotifyPlaylistSubscribers, websocket.PlaylistNeedRefresh)
	} else if vote.DownVote == true {
		resp = models.DownVoteSong(uint(songId), user, websocket.NotifyPlaylistSubscribers, websocket.PlaylistNeedRefresh)
	}
	u.Respond(w, resp)
}

func getSongId(params url.Values) (uint, error) {
	param1, ok := params["song_id"]
	if !ok || len(param1) < 1 {
		return 0, &net.ParseError{Type: "song_id", Text: param1[0]}
	}
	songId, err := strconv.ParseUint(param1[0], 10, 32)
	return uint(songId), err
}
