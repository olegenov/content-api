package models

import "github.com/jinzhu/gorm"

type Team struct {
	gorm.Model

	Name      string `json:"name" gorm:"not null"`
	CreatorID uint
	Creator   User
	Users     []User `json:"users" gorm:"many2many:team_users;"`
}
