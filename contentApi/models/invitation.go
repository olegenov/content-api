package models

import "github.com/jinzhu/gorm"

type Invitation struct {
	gorm.Model

	SenderUsername string `json:"sender_username" gorm:"not null"`
	ReceiverID     uint   `json:"receiver_id" gorm:"not null"`
	TeamID         uint   `json:"team_id" gorm:"not null"`
	Sender         User
	Receiver       User
	Team           Team
}
