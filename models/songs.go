package models

import (
	"fmt"
	u "github.com/PierreBougon/Bym-BackEnd/utils"
	"sort"

	"github.com/jinzhu/gorm"
)

type Song struct {
	gorm.Model
	Name       string `json:"name"`
	PlaylistId uint   `json:"playlist_id"`
	ExternalId string `json:"external_id"`
	VoteUp     int    `json:"vote_up"`
	VoteDown   int    `json:"vote_down"`
	Score      int    `json:"score"`
	Status     string `json:"status"`
	// We can add image + infos etc
}

// Not a model used to hold part of the song model
type Ranking struct {
	SongId   uint `json:"song_id"`
	VoteUp   int  `json:"vote_up"`
	VoteDown int  `json:"vote_down"`
	Score    int  `json:"score"`
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
	if err != nil /*|| playlist.UserId != user */ {
		return u.Message(false, "Invalid song, playlist may not be created"), false
	}
	if song.ExternalId == "" {
		return u.Message(false, "Invalid external id"), false
	}
	return u.Message(false, "Requirement passed"), true
}

func (song *Song) Create(user uint) map[string]interface{} {

	song.Status = "NONE"
	if resp, ok := song.Validate(user); !ok {
		return resp
	}

	GetDB().Create(song)

	if song.ID <= 0 {
		return u.Message(false, "Failed to create song, connection error.")
	}
	updatePlaylistSoundCount(user, song.PlaylistId, 1)
	response := u.Message(true, "song has been created")
	response["song"] = song
	return response
}

func GetSongs(playlist uint) []*Song {
	songs := make([]*Song, 0)
	err := GetDB().Table("songs").Where("playlist_id = ?", playlist).Order("score desc").Find(&songs).Error
	fmt.Println(err)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	//songs = pushFrontPlayingSong(songs)
	songs = pushFrontPlayedSongs(songs)
	return songs
}

func refreshPlaylistScoring(playlistId uint) {
	songs := make([]*Song, 0)
	err := GetDB().Table("songs").Where("playlist_id = ?", playlistId).Order("score desc").Find(&songs).Error
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := 0; i < len(songs); i++ {
		song := songs[i]
		song.Score += 10
		if song.Status == "PLAYED" {
			song.Score += song.VoteUp*10 - song.VoteDown*10
		} else if song.Status == "NONE" {
			song.Score += song.VoteUp*30 - song.VoteDown*30
		}
	}
	db.Save(songs)
}

func pushFrontPlayingSong(songs []*Song) []*Song {
	for i := 0; i < len(songs); i++ {
		if songs[i].Status == "PLAYING" {
			playingSong := songs[i]
			songs = append(songs[:i], songs[i:]...)
			songs = append([]*Song{playingSong}, songs...)
			return songs
		}
	}
	return nil
}

func pushFrontPlayedSongs(songs []*Song) []*Song {
	sort.Slice(songs, func(i, j int) bool {
		status := []string{
			"PLAYED",
			"PLAYING",
			"PAUSE",
			"NONE",
		}
		first := songs[i].Status
		second := songs[j].Status
		for k := 0; k < len(status); k++ {
			if first == second {
				return false
			}
			if status[k] == first {
				return true
			} else if status[k] == second {
				return false
			}
		}
		return true
	})
	return songs
}

func (song *Song) UpdateSong(user uint, songId uint, newSong *Song) map[string]interface{} {
	retSong := &Song{}
	err := db.First(&retSong, songId).Error
	playlist := &Playlist{}
	db.First(playlist, retSong.PlaylistId)
	if err != nil /*|| playlist.UserId != user*/ {
		return u.Message(false, "Invalid song")
	}
	//if (retSong.PlaylistId) TODO : very ownership
	if newSong.Name != "" {
		retSong.Name = newSong.Name
	}
	if newSong.Status != "" && isStatusValid(newSong.Status) {
		if newSong.Status == "STOP" {
			newSong.Status = "PLAYED"
			retSong.Score = -1
			go refreshPlaylistScoring(retSong.PlaylistId)
		}
		retSong.Status = newSong.Status
	}
	db.Save(&retSong)
	return u.Message(true, "Song successfully updated")
}

func isStatusValid(status string) bool {
	if status == "NONE" || status == "PLAYING" || status == "STOP" || status == "PAUSE" || status == "PLAYED" {
		return true
	}
	return false
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
	updatePlaylistSoundCount(user, song.PlaylistId, -1)
	return u.Message(true, "Song successfully deleted")
}

func GetSongsRanking(playlist uint) []*Ranking {

	songs := make([]*Song, 0)
	err := GetDB().Table("songs").Where("playlist_id = ?", playlist).Find(&songs).Order("score desc").Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	size := len(songs)
	songJson := Song{}
	ranks := make([]*Ranking, size)
	for j := 0; j < len(ranks); j++ {
		ranks[j] = &Ranking{}
	}
	for i := 0; i < size; i++ {
		songJson = *songs[i]
		ranks[i].SongId = songJson.ID
		ranks[i].VoteUp = songJson.VoteUp
		ranks[i].VoteDown = songJson.VoteDown
		ranks[i].Score = songJson.Score
	}
	return ranks
}

func GetSongRankingById(songid uint) *Ranking {

	song := Song{}
	err := GetDB().Table("songs").Where("id = ?", songid).Find(&song).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	rank := Ranking{}
	rank.SongId = song.ID
	rank.VoteUp = song.VoteUp
	rank.VoteDown = song.VoteDown
	rank.Score = song.Score

	return &rank
}

func RefreshSongVotes(songid uint) {
	//TODO : do not call this method directly create a thread to handle the refresh in another thread with a correct delay to not make 100 refresh/s for now just use it as it is
	//votes := make([]*Vote, 0)
	upVotes := 0
	downVotes := 0
	err1 := GetDB().Table("votes").Where("song_id = ? AND up_vote = ?", songid, true).Count(&upVotes).Error
	err2 := GetDB().Table("votes").Where("song_id = ? AND down_vote = ?", songid, true).Count(&downVotes).Error
	if err1 != nil || err2 != nil {
		return
	}
	song := Song{}
	err := db.Table("songs").First(&song, songid).Error
	if err != nil {
		return
	}
	song.VoteUp = upVotes
	song.VoteDown = downVotes
	song.Score = upVotes*100 - downVotes*100
	db.Save(&song)
}
