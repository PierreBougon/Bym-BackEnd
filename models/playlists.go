package models

import (
	"fmt"
	u "github.com/PierreBougon/Bym-BackEnd/utils"
	"github.com/jinzhu/gorm"
)

type Playlist struct {
	Model
	Name        string `json:"name"`
	UserId      uint   `json:"user_id"`
	SongsNumber int    `json:"songs_number"`
	Songs       []Song `gorm:"ForeignKey:PlaylistId"`
	Follower	[]*Account `gorm:"many2many:account_playlist;"`
	FollowerCount int 	`json:"follower_count"`
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
	playlist.SongsNumber = 0
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

func (playlist *Playlist) Join(user uint, playlistId uint) map[string]interface{} {
	account := &Account{}
	retPlaylist := &Playlist{}
	err := GetDB().Table("playlists").Where("id = ?", playlistId).Find(&retPlaylist).Error
	if err != nil {
		return u.Message(false, "This playlist does not exist")
	}
	if retPlaylist.UserId == user {
		return u.Message(true, "Author does not need to follow his playlist")
	}

	res := make([]*Account, 0)
	GetDB().Model(retPlaylist).Association("Follower").Find(&res)
	for _, follower := range res {
		if follower.ID == user {
			return u.Message(true, "User already joined the playlist")
		}
	}

	GetDB().Table("accounts").Where("id = ?", user).Find(&account)
	GetDB().Model(retPlaylist).Association("Follower").Append(account)
	GetDB().Model(retPlaylist).UpdateColumn("follower_count", gorm.Expr("follower_count + ?", 1))
	return u.Message(true, "User has joined the playlist")
}

func GetPlaylistById(u uint) *Playlist {
	retPlaylist := &Playlist{}
	GetDB().Preload("Songs").Table("playlists").Where("id = ?", u).First(retPlaylist)
	if retPlaylist.Name == "" {
		return nil
	}
	return retPlaylist
}

func GetPlaylistsByUser(user uint) []*Playlist {

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

func updatePlaylistSoundCount(user uint, playlistId uint, countModifier int) map[string]interface{} {
	retPlaylist := &Playlist{}
	err := db.Where(&Playlist{UserId: user}).First(&retPlaylist, playlistId).Error
	if err != nil {
		return u.Message(false, "Invalid playlist, you may not own this playlist")
	}
	retPlaylist.SongsNumber += countModifier
	db.Save(&retPlaylist)
	return u.Message(true, "Playlist successfully updated")
}

func (playlist *Playlist) DeletePlaylist(user uint, playlistId uint) map[string]interface{} {
	retPlaylist := &Playlist{}
	err := db.Where(&Playlist{UserId: user}).First(&retPlaylist, playlistId).Error
	if err != nil {
		return u.Message(false, "Invalid playlist, you may not own this playlist")
	}
	GetDB().Model(retPlaylist).Association("Follower").Clear()
	db.Delete(&retPlaylist)
	return u.Message(true, "Playlist successfully deleted")
}
