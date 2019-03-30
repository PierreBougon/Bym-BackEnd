package models

import (
	u "github.com/PierreBougon/Bym-BackEnd/utils"

	"fmt"

	"github.com/jinzhu/gorm"
)

type Playlist struct {
	gorm.Model
	Name   string `json:"name"`
	UserId uint   `json:"user_id"`
	Songs  []Song
}

func (playlist *Playlist) Validate() (map[string]interface{}, bool) {
	if len(playlist.Name) < 3 {
		return u.Message(false, "Playlist name must be at least 3 characters long"), false
	}
	if playlist.UserId <= 0 {
		return u.Message(false, "Invalid user"), false
	}
	return u.Message(false, "Requirement passed"), true
}

func (playlist *Playlist) Create(user uint) map[string]interface{} {

	playlist.UserId = user
	if resp, ok := playlist.Validate(); !ok {
		return resp
	}

	GetDB().Create(playlist)

	if playlist.ID <= 0 {
		return u.Message(false, "Failed to create playlist, connection error.")
	}

	response := u.Message(true, "Playlist has been created")
	response["playlist"] = playlist
	return response
}

func GetPlaylistById(u uint) *Playlist {
	playlist := &Playlist{}
	GetDB().Table("playlists").Where("id = ?", u).First(playlist)
	if playlist.Name == "" {
		return nil
	}
	return playlist
}

func GetPlaylistByUser(user uint) []*Playlist {

	playlists := make([]*Playlist, 0)
	err := GetDB().Table("playlists").Where("user_id = ?", user).Find(&playlists).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return playlists
}
