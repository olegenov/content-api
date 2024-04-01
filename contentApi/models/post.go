package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Post struct {
	gorm.Model

	Title       string `json:"title"`
	AssignID    uint
	Assign      User `json:"assign" gorm:"foreignKey:AssignID"`
	PublishDate time.Time
	Deadline    time.Time
	ProjectID   uint
	Project     Project
	Tags        []Tag  `json:"tags" gorm:"many2many:post_tags;"`
	Content     string `json:"content" gorm:"type:text"`
}
