package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Post struct {
	gorm.Model

	Title       string `json:"title" gorm:"not null"`
	PublishDate time.Time
	Deadline    time.Time
	ProjectID   uint
	Tags        []Tag  `json:"tags" gorm:"many2many:post_tags;"`
	Content     string `json:"content" gorm:"type:text"`
}
