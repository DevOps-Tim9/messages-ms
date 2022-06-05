package entity

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	From           uint   `gorm:"not null;default:null"`
	To             uint   `gorm:"not null;default:null"`
	Text           string `gorm:"not null;default:null"`
	ConversationId uint   `gorm:"not null;default:null"`

	Tbl string `gorm:"-"`
}
