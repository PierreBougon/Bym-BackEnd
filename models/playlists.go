package models

import (
	"fmt"
	u "github.com/PierreBougon/Bym-BackEnd/utils"

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

func GetPlaylist(u uint) *Playlist {
	retPlaylist := &Playlist{}
	GetDB().Table("playlists").Where("id = ?", u).First(retPlaylist)
	if retPlaylist.Name == "" {
		return nil
	}
	return retPlaylist
}

func GetPlaylists(user uint) []*Playlist {

	playlists := make([]*Playlist, 0)
	err := GetDB().Table("playlists").Where("user_id = ?", user).Find(&playlists).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return playlists
}

func (playlist *Playlist) UpdatePlaylist(user uint, playlistId uint, newPlaylist *Playlist) map[string]interface{} {
	retPlaylist := &Playlist{}
	err := db.Where(&Playlist{UserId: user}).First(&retPlaylist, playlistId).Error
	if err != nil {
		return u.Message(false, "Invalid playlist, you may not own this playlist")
	}
	retPlaylist.Name = newPlaylist.Name
	db.Save(&retPlaylist)
	return u.Message(true, "Playlist successfully updated")
}

func (playlist *Playlist) DeletePlaylist(user uint, playlistId uint) map[string]interface{} {
	retPlaylist := &Playlist{}
	err := db.Where(&Playlist{UserId: user}).First(&retPlaylist, playlistId).Error
	if err != nil {
		return u.Message(false, "Invalid playlist, you may not own this playlist")
	}
	db.Delete(&retPlaylist)
	return u.Message(true, "Playlist successfully deleted")
}
