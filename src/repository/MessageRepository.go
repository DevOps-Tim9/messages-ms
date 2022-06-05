package repository

import (
	"messages-ms/src/entity"

	"gorm.io/gorm"
)

type IMessageRepository interface {
	Create(entity.Message) (entity.Message, error)
	GetMesssagesByConversation(uint) []*entity.Message
}

type MessageRepository struct {
	Database *gorm.DB
}

func (r MessageRepository) Create(message entity.Message) (entity.Message, error) {
	err := r.Database.Save(&message).Error

	return message, err
}

func (r MessageRepository) GetMesssagesByConversation(conversationId uint) []*entity.Message {
	var messages = []*entity.Message{}

	r.Database.Preload("Conversation").Where("conversation_id = ?", conversationId).Order("created_at DESC").Find(&messages)

	return messages
}
