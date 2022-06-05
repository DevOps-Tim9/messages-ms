package repository

import (
	"messages-ms/src/entity"
	"strconv"

	"gorm.io/gorm"
)

type IConversationRepository interface {
	Create(entity.Conversation) (entity.Conversation, error)
	GetConversationByUsers(uint, uint) (*entity.Conversation, error)
	GetConversationsByUser(uint) []*entity.Conversation
}

type ConversationRepository struct {
	Database *gorm.DB
}

func (r ConversationRepository) Create(conversation entity.Conversation) (entity.Conversation, error) {
	err := r.Database.Save(&conversation).Error

	return conversation, err
}

func (r ConversationRepository) GetConversationByUsers(user1 uint, user2 uint) (*entity.Conversation, error) {
	var conversation = entity.Conversation{}

	err := r.Database.
		Where("(user1 = " + strconv.FormatUint(uint64(user1), 10) + " AND user2 = " + strconv.FormatUint(uint64(user2), 10) + ") OR (user1 = " + strconv.FormatUint(uint64(user2), 10) + " AND user2 = " + strconv.FormatUint(uint64(user1), 10) + ")").First(&conversation).Error

	return &conversation, err
}

func (r ConversationRepository) GetConversationsByUser(userId uint) []*entity.Conversation {
	var conversations = []*entity.Conversation{}

	r.Database.Preload("LastMessage").Where("user1 = " + strconv.FormatUint(uint64(userId), 10) + " OR user2 = " + strconv.FormatUint(uint64(userId), 10)).Order("last_message_id DESC").Find(&conversations)

	return conversations
}
