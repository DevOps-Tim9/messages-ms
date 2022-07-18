package repository

import (
	"context"
	"errors"
	"messages-ms/src/entity"
	"time"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageRepositoryMock struct {
	mock.Mock
}

func (m MessageRepositoryMock) Create(message entity.Message, ctx context.Context) (entity.Message, error) {
	if message.From == 5 {
		return entity.Message{}, errors.New("")
	}

	message.ID = primitive.NewObjectID()

	return message, nil
}

func (m MessageRepositoryMock) Update(message entity.Message, ctx context.Context) error {
	if message.ID == primitive.NilObjectID {
		return errors.New("")
	}

	return nil
}

func (m MessageRepositoryMock) GetMesssagesByConversation(conversationId string, ctx context.Context) []entity.Message {
	id, err := primitive.ObjectIDFromHex(conversationId)

	if err != nil {
		return make([]entity.Message, 0)
	}

	return []entity.Message{
		{
			ID:             primitive.NewObjectID(),
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
			From:           1,
			To:             2,
			ConversationId: id,
			Text:           "Text",
		},
	}
}
