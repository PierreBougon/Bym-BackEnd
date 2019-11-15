package models

import (
	"fmt"
	u "github.com/PierreBougon/Bym-BackEnd/utils"
	"github.com/jinzhu/gorm"
)

type Playlist struct {
	Model
	Name          string                  `json:"name"`
	UserId        uint                    `json:"user_id"`
	SongsNumber   int                     `json:"songs_number"`
	Songs         []Song                  `gorm:"ForeignKey:PlaylistId"`
	Follower      []*Account              `gorm:"many2many:account_playlist;"`
	FollowerCount int                     `json:"follower_count"`
	Acl           []PlaylistAccessControl `gorm:"ForeignKey:PlaylistId"`
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

// TODO(variable POST optionel pour choisir le ROLE de l'user (BYMER ou FOllOWER), BYMER by default)
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
	// TODO(do the append only if Bymer ROLE is asked in the post parameter)
	GetDB().Model(retPlaylist).Association("Acl").Append(&PlaylistAccessControl{
		UserId:     user,
		PlaylistId: retPlaylist.ID,
		RoleId:     ROLE_BYMER,
	})
	return u.Message(true, "User has joined the playlist")
}

/*
TODO(Cas où l'utilisateur dont l'acl est changé n'appartient pas à la playlist ????)
TODO(Cas où le front ne peut envoyer qu'un email et pas un id)
TODO(Roles sont rentré à la main dans la db, dans l'ordre croissant d'authorité (1 > 4))
TODO(Le front peut-il directement utilisé ces 'valeurs' d'autorité en dure, sans avoir besoin de les fetch dans la db ?)
*/
func ChangeAclOnPlaylist(user uint, userToPromote uint, playlistId uint, role uint) map[string]interface{} {
	var (
		userAcl     uint
		retPlaylist Playlist
		oldAcl      PlaylistAccessControl
	)

	db.Table("playlist_access_controls").
		Where(&PlaylistAccessControl{UserId: user, PlaylistId: playlistId}).Pluck("role", userAcl)
	db.Table("playlists").Where("id = ?", playlistId).Find(&retPlaylist)
	if userAcl != ROLE_ADMIN && retPlaylist.UserId != user {
		return u.Message(false, "User has no right upon the members of this playlist")
	}
	// At this point, if userAcl == 0 then he is Author
	exist := !db.Table("playlist_access_controls").
		Where("user_id = ? AND playlist_id = ?", userToPromote, playlistId).Find(&oldAcl).
		RecordNotFound()
	if oldAcl.RoleId == role {
		return u.Message(true, "This user already have this role")
	}
	if oldAcl.RoleId == ROLE_ADMIN && role > ROLE_ADMIN && userAcl != 0 {
		return u.Message(false, "To demote an Admin needs Author rights")
	}
	if exist {
		db.Table("playlist_access_controls").
			Where("playlist_id = ? AND user_id = ?", playlistId, userToPromote).
			UpdateColumn("role_id", role)
	} else {
		return u.Message(false, "Can not find this user")
	}
	return u.Message(true, "New role successfully given")
}

func (playlist *Playlist) LeavePlaylist(user uint, playlistId uint) map[string]interface{} {
	retPlaylist := &Playlist{}
	err := db.Where("id = ?", playlistId).First(&retPlaylist, playlistId).Error
	if err != nil {
		return u.Message(false, "Playlist does not exist")
	}
	accounts := make([]*Account, 0)
	GetDB().Model(retPlaylist).Association("Follower").Find(&accounts)
	isFollowed := false
	for _, account := range accounts {
		if account.ID == user {
			isFollowed = true
		}
	}
	if !isFollowed {
		return u.Message(false, "Playlist is not followed by user")
	}
	account := &Account{}
	GetDB().Table("accounts").Find(&account, "id = ?", user)
	GetDB().Model(retPlaylist).Association("Follower").Delete(&account)
	GetDB().Model(retPlaylist).UpdateColumn("follower_count", gorm.Expr("follower_count - ?", 1))
	return u.Message(true, "Playlist successfully left")
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
	db.Model(retPlaylist).Association("Follower").Clear()
	db.Model(retPlaylist).Association("Acl").Clear()
	db.Delete(retPlaylist)
	return u.Message(true, "Playlist successfully deleted")
}

func GetRole(user uint, playlistId uint) (*Role, string)  {
	retPlaylist := &Playlist{}
	err := GetDB().Table("playlists").Where("id = ?", playlistId).First(&retPlaylist).Error
	if err != nil {
		return nil, "Invalid playlist, it does not exist"
	}
	if retPlaylist.UserId == user {
		return &Role{Name: RoleName[0], ID: 0}, ""
	}

	acl := &PlaylistAccessControl{};
	notFound := db.Table("playlist_access_controls").Where(PlaylistAccessControl{
		UserId:     user,
		PlaylistId: playlistId,
	}).First(acl).RecordNotFound()
	if notFound {
		return &Role{Name: RoleName[ROLE_VISITOR], ID: ROLE_VISITOR}, ""
	}
	return &Role{ID: acl.RoleId, Name: RoleName[acl.RoleId]}, ""
}
