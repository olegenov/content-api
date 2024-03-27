package models

import "github.com/jinzhu/gorm"

type Project struct {
	gorm.Model

	Name      string `json:"name" gorm:"not null"`
	CreatorID uint
	Creator   User
	Posts     []Post `json:"posts" gorm:"foreignKey:ProjectID,constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	TeamID    uint
	Team      Team
}
