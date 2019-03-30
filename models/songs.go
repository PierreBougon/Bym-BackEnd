package models

import (
	u "github.com/PierreBougon/Bym-BackEnd/utils"

	"fmt"

	"github.com/jinzhu/gorm"
)

type Song struct {
	gorm.Model
	Name       string `json:"name"`
	PlaylistId uint   `json:"playlist_id"`
	ExternalId string `json:"external_id"`
	// We can add image + infos etc
}

func (song *Song) Validate() (map[string]interface{}, bool) {
	if len(song.Name) < 1 {
		return u.Message(false, "Song name must be at least 1 characters long"), false
	}
	if song.PlaylistId == 0 {
		return u.Message(false, "Invalid playlist"), false
	}
	if song.ExternalId == "" {
		return u.Message(false, "Invalid external id"), false
	}
	return u.Message(false, "Requirement passed"), true
}

func (song *Song) Create() map[string]interface{} {

	if resp, ok := song.Validate(); !ok {
		// fmt.Println(resp, ok)
		return resp
	}

	GetDB().Create(song)

	if song.ID <= 0 {
		return u.Message(false, "Failed to create song, connection error.")
	}

	response := u.Message(true, "song has been created")
	response["song"] = song
	return response
}

func GetSongs(playlist uint) []*Song {

	songs := make([]*Song, 0)
	err := GetDB().Table("songs").Where("playlist_id = ?", playlist).Find(&songs).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return songs
}
