package models

import "github.com/jinzhu/gorm"

type Tag struct {
	gorm.Model

	Name  string `gorm:"unique;not null"`
	Color string `gorm:"not null"`
	Posts []Post `gorm:"gorm:many2many:post_tags"`
}
