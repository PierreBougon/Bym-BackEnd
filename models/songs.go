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

func (song *Song) Validate(user uint) (map[string]interface{}, bool) {
	if len(song.Name) < 1 {
		return u.Message(false, "Song name must be at least 1 characters long"), false
	}
	if song.PlaylistId == 0 {
		return u.Message(false, "Invalid playlist"), false
	}
	playlist := &Playlist{}
	err := db.First(playlist, song.PlaylistId).Error
	if err != nil || playlist.UserId != user {
		return u.Message(false, "Invalid song, you may not own this playlist or playlist doesn't exist"), false
	}
	if song.ExternalId == "" {
		return u.Message(false, "Invalid external id"), false
	}
	return u.Message(false, "Requirement passed"), true
}

func (song *Song) Create(user uint) map[string]interface{} {

	if resp, ok := song.Validate(user); !ok {
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

func (song *Song) UpdateSong(user uint, songId uint, newSong *Song) map[string]interface{} {
	retSong := &Song{}
	err := db.First(&retSong, songId).Error
	playlist := &Playlist{}
	db.First(playlist, retSong.PlaylistId)
	if err != nil || playlist.UserId != user {
		return u.Message(false, "Invalid song, you may not own this playlist")
	}
	//if (retSong.PlaylistId) TODO : very ownership
	retSong.Name = newSong.Name
	db.Save(&retSong)
	return u.Message(true, "Song successfully updated")
}

func (song *Song) DeleteSong(user uint, songId uint) map[string]interface{} {
	retSong := &Song{}
	err := db.First(&retSong, songId).Error
	playlist := &Playlist{}
	db.First(playlist, retSong.PlaylistId)
	if err != nil || playlist.UserId != user {
		return u.Message(false, "Invalid song, you may not own this playlist")
	}
	db.Delete(&retSong)
	return u.Message(true, "Song successfully deleted")
}
