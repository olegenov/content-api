package models

import "github.com/jinzhu/gorm"

type Project struct {
	gorm.Model

	Name      string `gorm:"not null"`
	CreatorID uint
	Posts     []Post `gorm:"foreignKey:ProjectID,constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
