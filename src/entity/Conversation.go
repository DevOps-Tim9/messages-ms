package entity

import "gorm.io/gorm"

type Conversation struct {
	gorm.Model
	User1         uint `gorm:"not null;default:null"`
	User2         uint `gorm:"not null;default:null"`
	LastMessageId uint
	LastMessage   Message

	Tbl string `gorm:"-"`
}
