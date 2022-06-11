package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName   string `gorm:"not null" json:"UserName"`
	UserPwd    string `gorm:"not null" json:"UserPwd"`
	University string `json:"University"`
}

type Message struct {
	gorm.Model
	ConversationID int
	Message        string `gorm:"not null"`
	SenderId       int    `gorm:"not null"`
	ReceiverId     int    `gorm:"not null"`
}

type Conversation struct {
	gorm.Model
	User1Id  int
	User2Id  int
	User1    User
	User2    User
	Messages []Message
}
