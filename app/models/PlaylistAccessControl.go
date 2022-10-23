package models

import (
	"fmt"
	"strconv"
	"time"
)

type PlaylistAccessControl struct {
	CreatedAt  time.Time
	UpdatedAt  time.Time
	UserId     uint `gorm:"unique_index:idx_user_playlist" json:"user_id"`
	PlaylistId uint `gorm:"unique_index:idx_user_playlist" json:"playlist_id"`
	RoleId     uint `json:"role_id"`
}

func checkRight(user uint, playlist uint, neededRole uint) bool {
	acl := &PlaylistAccessControl{}
	p := &Playlist{}
	fmt.Println("Check right on playlist " + fmt.Sprintf("%d role needed is %d", playlist, neededRole))
	db.Table("playlists").Where("id = ?", playlist).First(p)
	if user == p.UserId {
		fmt.Println("User is authorized to do the action because he is the author.")
		return true
	}

	notFound := db.Table("playlist_access_controls").Where(PlaylistAccessControl{
		UserId:     user,
		PlaylistId: playlist,
	}).First(acl).RecordNotFound()
	if notFound {
		fmt.Println("User does not have any right on this playlist.")
		return false
	}
	fmt.Println("Current role on this playlist : " + fmt.Sprint(acl.RoleId) + " Result : " + strconv.FormatBool(acl.RoleId <= neededRole))
	return acl.RoleId <= neededRole
}
