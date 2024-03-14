package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Post struct {
	gorm.Model

	Title       string `gorm:"not null"`
	PublishDate time.Time
	Deadline    time.Time
	ProjectID   uint
	Tags        []Tag  `gorm:"many2many:post_tags;"`
	Content     string `gorm:"type:text"`
}
