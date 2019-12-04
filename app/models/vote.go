package models

import (
	u "github.com/PierreBougon/Bym-BackEnd/app/utils"
)

type Vote struct {
	Model
	UpVote   bool `json:"up_vote"`
	DownVote bool `json:"down_vote"`
	UserId   uint `json:"user_id"`
	SongId   uint `json:"song_id"`
}

func GetPersonalVoteBySongId(songid uint, user uint) *Vote {
	vote := Vote{}
	err := GetDB().Table("votes").Where("song_id = ? AND user_id = ?", songid, user).Find(&vote).Error
	if err != nil {
		//fmt.Println(err)
		return nil
	}
	return &vote
}

func GetVotesBySongId(songid uint) []*Vote {
	votes := make([]*Vote, 0)
	err := GetDB().Table("votes").Where("song_id = ?", songid).Find(&votes).Error
	if err != nil {
		return nil
	}
	return votes
}

func updateVote(songid uint, user uint, upVote bool, notifyOnUpdate func(userId uint, playlistId uint, message string), messageOnUpdate func(playlistId uint, userId uint) string) map[string]interface{} {
	err := GetDB().Table("songs").Find(&Song{}, "id = ?", songid).Error
	if err != nil {
		return u.Message(false, "Request failed, connection error or songId does not exist")
	}

	vote := Vote{}
	res := GetDB().Table("votes").
		Find(&vote, "song_id = ? AND user_id = ?", songid, user)
	errNbr := len(res.GetErrors())
	notFound := res.RecordNotFound()
	// Database failure: Only one error happened which is not RecordNotFound or other error(s) happened
	if (errNbr > 1 && notFound) || (errNbr > 0 && !notFound) {
		return u.Message(false, "Request failed, connection error")
		// If Vote did not exist, fill the data of the new one
	} else if res.RecordNotFound() {
		vote.UserId = user
		vote.SongId = songid
	} else {
		if vote.UpVote == upVote {
			return u.Message(true, "This request has performed no action")
		}
	}
	// Update value of the up/down vote
	vote.UpVote = upVote
	vote.DownVote = !upVote
	db.Save(&vote)
	go RefreshSongVotes(user, songid, notifyOnUpdate, messageOnUpdate)
	return u.Message(true, "Song successfully up voted !")
}

func UpVoteSong(songid uint, user uint, notifyOnUpdate func(userId uint, playlistId uint, message string), messageOnUpdate func(playlistId uint, userId uint) string) map[string]interface{} {
	return updateVote(songid, user, true, notifyOnUpdate, messageOnUpdate)
}

func DownVoteSong(songid uint, user uint, notifyOnUpdate func(userId uint, playlistId uint, message string), messageOnUpdate func(playlistId uint, userId uint) string) map[string]interface{} {
	return updateVote(songid, user, false, notifyOnUpdate, messageOnUpdate)
}
