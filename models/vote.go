package models

import (
	"fmt"
	u "github.com/PierreBougon/Bym-BackEnd/utils"
	"github.com/jinzhu/gorm"
)

type Vote struct {
	gorm.Model
	UpVote   bool `json:"up_vote"`
	DownVote bool `json:"down_vote"`
	UserId   uint `json:"user_id"`
	SongId   uint `json:"song_id"`
}

func GetPersonalVoteBySongId(songid uint, user uint) *Vote {
	vote := Vote{}
	err := GetDB().Table("votes").Where("song_id = ? AND user_id = ?", songid, user).Find(&vote).Error
	if err != nil {
		fmt.Println(err)
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

func UpVoteSong(songid uint, user uint) map[string]interface{} {
	flag := false
	vote := Vote{}
	err := GetDB().Table("votes").Where("song_id = ? AND user_id = ?", songid, user).Find(&vote).Error
	if err != nil {
		flag = true
		vote.UserId = user
		vote.SongId = songid
	}

	if vote.UpVote == true {
		return u.Message(true, "This request has performed no action")
	}
	if vote.DownVote == true {
		vote.DownVote = false
	}
	if vote.UpVote == false {
		vote.UpVote = true
		if flag {
			db.Create(&vote)
		} else {
			db.Save(&vote)
		}
		RefreshSongVotes(songid)
	}
	return u.Message(true, "Song successfully up voted !")
}

func DownVoteSong(songid uint, user uint) map[string]interface{} {
	flag := false
	vote := Vote{}
	err := GetDB().Table("votes").Where("song_id = ? AND user_id = ?", songid, user).Find(&vote).Error
	if err != nil {
		flag = true
		vote.UserId = user
		vote.SongId = songid
	}

	if vote.DownVote == true {
		return u.Message(true, "This request has performed no action")
	}
	if vote.UpVote == true {
		vote.UpVote = false
	}
	if vote.DownVote == false {
		vote.DownVote = true
		if flag {
			db.Create(&vote)
		} else {
			db.Save(&vote)
		}
		RefreshSongVotes(songid)
	}
	return u.Message(true, "Song successfully down voted !")
}
