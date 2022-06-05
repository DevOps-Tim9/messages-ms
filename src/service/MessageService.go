package service

import (
	"messages-ms/src/dto"
	"messages-ms/src/entity"
	"messages-ms/src/repository"
)

type IMessageService interface {
	CreateNewMessage(dto.MessageDto) (entity.Message, error)
	GetMesssagesByConversation(uint) []*entity.Message
	GetConversationsByUser(uint) []*entity.Conversation
	GetConversationByUsers(uint, uint) (*entity.Conversation, error)
}

type MessageService struct {
	MessageRepository      repository.IMessageRepository
	ConversationRepository repository.IConversationRepository
}

func (s MessageService) GetMesssagesByConversation(id uint) []*entity.Message {
	return s.MessageRepository.GetMesssagesByConversation(id)
}

func (s MessageService) GetConversationsByUser(id uint) []*entity.Conversation {
	return s.ConversationRepository.GetConversationsByUser(id)
}

func (s MessageService) GetConversationByUsers(user1 uint, user2 uint) (*entity.Conversation, error) {
	return s.ConversationRepository.GetConversationByUsers(user1, user2)
}

func (s MessageService) CreateNewMessage(dto dto.MessageDto) (entity.Message, error) {
	conversation, err := s.ConversationRepository.GetConversationByUsers(dto.From, dto.To)

	if err != nil {
		conversation = &entity.Conversation{
			User1: dto.From,
			User2: dto.To,
		}

		newConversation, _ := s.ConversationRepository.Create(*conversation)

		newMessage, _ := s.MessageRepository.Create(entity.Message{
			From:           dto.From,
			To:             dto.To,
			Text:           dto.Text,
			ConversationId: newConversation.ID,
		})

		newConversation.LastMessageId = newMessage.ID

		s.ConversationRepository.Create(newConversation)

		return newMessage, nil
	} else {
		newMessage, err := s.MessageRepository.Create(entity.Message{
			From:           dto.From,
			To:             dto.To,
			Text:           dto.Text,
			ConversationId: conversation.ID,
		})

		conversation.LastMessageId = newMessage.ID

		s.ConversationRepository.Create(*conversation)

		return newMessage, err
	}
}
