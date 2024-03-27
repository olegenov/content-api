package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model

	Email     string    `json:"email" gorm:"not null"`
	Firstname string    `json:"firstname"`
	Surname   string    `json:"surname"`
	Username  string    `json:"username" gorm:"unique;not null"`
	Password  string    `json:"password" gorm:"not null"`
	Projects  []Project `gorm:"foreignKey:CreatorID"`
	Role      string
	Teams     []Team `gorm:"gorm:many2many:team_users"`
}
