package models

import "time"

type PlaylistAccessControl struct {
	CreatedAt 	time.Time
	UpdatedAt 	time.Time
	UserId		uint	`gorm:"unique_index:idx_user_playlist",json:"user_id"`
	PlaylistId 	uint	`gorm:"unique_index:idx_user_playlist",json:"playlist_id"`
	RoleId      uint	`json:"role_id"`
}

func checkRight(user uint, playlist uint, neededRole uint) bool {
	acl := &PlaylistAccessControl{}
	p := &Playlist{}
	db.Table("playlists").Where("id = ?", playlist).First(p)
	if user == p.UserId {
		return true
	}

	notFound := db.Table("playlist_access_control").Where(PlaylistAccessControl{
		UserId:     user,
		PlaylistId: playlist,
	}).First(acl).RecordNotFound()
	if notFound || acl.RoleId > neededRole {
		return false
	}
	return true
}