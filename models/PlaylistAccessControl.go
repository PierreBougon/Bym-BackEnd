package models

import "time"

type PlaylistAccessControl struct {
	CreatedAt 	time.Time
	UpdatedAt 	time.Time
	UserId		uint	`gorm:"unique_index:idx_user_playlist",json:"user_id"`
	PlaylistId 	uint	`gorm:"unique_index:idx_user_playlist",json:"playlist_id"`
	RoleId      uint	`json:"role_id"`
}
