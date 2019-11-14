package models

type Role struct {
	ID      uint   `gorm:"primary_key"`
	Name	string `json:"name"`
	Acl		[]PlaylistAccessControl `gorm:"ForeignKey:RoleId"`
	// We can add image + infos etc
}

// const ROLE_DEV =
const ROLE_ADMIN = 1
const ROLE_BYMER = 2
const ROLE_FOLLOWER  = 3
// const ROLE_VISITOR  =

