package models

type AccountPlaylist struct {
	AccountId  uint `json:"account_id"`
	PlaylistId uint `json:"playlist_id"`
}

func GetFollowers(playlistId uint) []*AccountPlaylist {
	accPlay := make([]*AccountPlaylist, 0)
	err := GetDB().Table("account_playlist").Where("playlist_id = ?", playlistId).Find(&accPlay).Error
	if err != nil {
		return nil
	}

	return accPlay
}

func CleanFollowersPersonalData(followers []*Account) {
	for _, follower := range followers {
		follower.TokenVersion = 0
		follower.Email = ""
		follower.Password = ""
	}
}
