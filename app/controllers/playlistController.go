package controllers

import (
	"fmt"
	"github.com/PierreBougon/Bym-BackEnd/app/models"
	u "github.com/PierreBougon/Bym-BackEnd/app/utils"
	"github.com/PierreBougon/Bym-BackEnd/websocket"

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
	fmt.Println("get playlists")
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
	filter := models.PlaylistFilter{
		ShowSongs:    true,
		ShowFollower: false,
		ShowAcl:      false,
	}
	getBoolVal := func(val []string, b *bool) bool {
		if len(val) >= 1 {
			ret, err := strconv.ParseBool(val[0])
			if err == nil {
				*b = ret
				return true
			}
		}
		return false
	}
	vals := r.URL.Query()
	for name, val := range vals {
		goodParam := false
		switch name {
		case "showSongs":
			goodParam = getBoolVal(val, &filter.ShowSongs)
		case "showFollower":
			goodParam = getBoolVal(val, &filter.ShowFollower)
		case "showAcl":
			goodParam = getBoolVal(val, &filter.ShowAcl)
		}
		if !goodParam {
			w.WriteHeader(http.StatusBadRequest)
			u.Respond(w, u.Message(false, "Invalid request, wrong value given to "+name))
			return
		}
	}
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		u.RespondBadRequest(w)
		return
	}

	data := models.GetPlaylistById(uint(id), &filter)
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
	} else {
		websocket.NotifyPlaylistSubscribers(user, uint(id), websocket.PlaylistNeedRefresh(uint(id), user))
	}
	u.Respond(w, resp)
}

var LeavePlaylist = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 32)
	if err != nil {
		u.RespondBadRequest(w)
		return
	}
	user := r.Context().Value("user").(uint)
	resp := (&models.Playlist{}).LeavePlaylist(user, uint(id))
	if resp["status"] == false {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		websocket.NotifyPlaylistSubscribers(user, uint(id), websocket.PlaylistNeedRefresh(uint(id), user))
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
	messageOnDelete := websocket.PlaylistDeleted(uint(id), user)
	resp := (&models.Playlist{}).DeletePlaylist(user, uint(id), websocket.NotifyPlaylistSubscribers, messageOnDelete)
	if resp["status"] == false {
		w.WriteHeader(http.StatusBadRequest)
	}
	u.Respond(w, resp)
}

var JoinPlaylist = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 32)
	if err != nil {
		u.RespondBadRequest(w)
		return
	}
	user := r.Context().Value("user").(uint)
	resp := (&models.Playlist{}).Join(user, uint(id))
	if resp["status"] == false {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		websocket.NotifyPlaylistSubscribers(user, uint(id), websocket.PlaylistNeedRefresh(uint(id), user))
	}
	u.Respond(w, resp)
}

var ChangeAclOnPlaylist = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 32)
	if err != nil {
		u.RespondBadRequest(w)
		return
	}

	type UserAcl struct {
		User *uint `json:"user"`
		Role *uint `json:"role"`
	}
	userAcl := &UserAcl{}
	err = json.NewDecoder(r.Body).Decode(userAcl)
	if err != nil || userAcl.Role == nil || userAcl.User == nil {
		u.RespondBadRequest(w)
		return
	}
	user := r.Context().Value("user").(uint)
	resp := models.ChangeAclOnPlaylist(user, *userAcl.User, uint(id), *userAcl.Role)
	if resp["status"] == false {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		websocket.NotifyPlaylistSubscribers(user, uint(id), websocket.PlaylistNeedRefresh(uint(id), user))
	}
	u.Respond(w, resp)
}

var GetPlaylistRole = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 32)
	if err != nil {
		u.RespondBadRequest(w)
		return
	}
	user := r.Context().Value("user").(uint)
	data, errMsg := models.GetRole(user, uint(id))
	var resp map[string]interface{}
	if errMsg != "" {
		resp = u.Message(false, errMsg)
	} else {
		resp = u.Message(true, "success")
	}
	resp["role"] = data
	u.Respond(w, resp)
}
