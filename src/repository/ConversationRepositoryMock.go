package repository

import (
	"context"
	"errors"
	"messages-ms/src/entity"
	"time"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConversationRepositoryMock struct {
	mock.Mock
}

func (c ConversationRepositoryMock) Create(conversation entity.Conversation, ctx context.Context) (entity.Conversation, error) {
	if conversation.User1 == 5 {
		return entity.Conversation{}, errors.New("")
	}

	conversation.ID = primitive.NewObjectID()

	return conversation, nil
}

func (c ConversationRepositoryMock) Update(conversation entity.Conversation, ctx context.Context) error {
	if conversation.ID == primitive.NilObjectID {
		return errors.New("")
	}

	return nil
}

func (c ConversationRepositoryMock) GetConversationByUsers(user1 uint, user2 uint, ctx context.Context) (entity.Conversation, error) {
	if user1 == 1 && user2 == 2 {
		return entity.Conversation{}, errors.New("")
	}

	conversationId := primitive.NewObjectID()

	return entity.Conversation{ID: conversationId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		User1:     2,
		User2:     3,
		LastMessage: entity.Message{
			ID:             primitive.NewObjectID(),
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
			From:           2,
			To:             3,
			ConversationId: conversationId,
			Text:           "Text",
		},
	}, nil
}

func (c ConversationRepositoryMock) GetConversationsByUser(id uint, ctx context.Context) []entity.Conversation {
	if id == 1 {
		return make([]entity.Conversation, 0)
	}

	conversationId := primitive.NewObjectID()

	return []entity.Conversation{
		{ID: conversationId,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			User1:     2,
			User2:     3,
			LastMessage: entity.Message{
				ID:             primitive.NewObjectID(),
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
				From:           2,
				To:             3,
				ConversationId: conversationId,
				Text:           "Text",
			},
		},
	}
}
