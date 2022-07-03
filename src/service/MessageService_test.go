package service

import (
	"context"
	"messages-ms/src/dto"
	"messages-ms/src/repository"
	"messages-ms/src/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageServiceUnitTestSuite struct {
	suite.Suite
	messageRepositoryMock      *repository.MessageRepositoryMock
	conversationRepositoryMock *repository.ConversationRepositoryMock
	service                    MessageService
}

func TestMessageServiceUnitTestSuite(t *testing.T) {
	suite.Run(t, new(MessageServiceUnitTestSuite))
}

func (suite *MessageServiceUnitTestSuite) SetupSuite() {
	suite.messageRepositoryMock = new(repository.MessageRepositoryMock)
	suite.conversationRepositoryMock = new(repository.ConversationRepositoryMock)

	suite.service = MessageService{
		MessageRepository:      suite.messageRepositoryMock,
		ConversationRepository: suite.conversationRepositoryMock,
		Logger:                 utils.Logger(),
	}
}

func (suite *MessageServiceUnitTestSuite) TestNewMessageService() {
	assert.NotNil(suite.T(), suite.service, "Service is nil")
}

func (suite *MessageServiceUnitTestSuite) TestMessageService_GetMesssagesByConversation_ReturnsEmptyList() {
	messages := suite.service.GetMesssagesByConversation("1", context.TODO())

	assert.NotNil(suite.T(), messages, "Messages are nil")
	assert.Equal(suite.T(), 0, len(messages), "Length of messages is not 0")
}

func (suite *MessageServiceUnitTestSuite) TestMessageService_GetMesssagesByConversation_ReturnsListOfMessages() {
	id := primitive.NewObjectID()

	messages := suite.service.GetMesssagesByConversation(id.Hex(), context.TODO())

	assert.NotNil(suite.T(), messages, "Messages are nil")
	assert.Equal(suite.T(), 1, len(messages), "Length of messages is not 1")
}

func (suite *MessageServiceUnitTestSuite) TestMessageService_GetConversationsByUser_ReturnsEmptyList() {
	conversations := suite.service.GetConversationsByUser(1, context.TODO())

	assert.NotNil(suite.T(), conversations, "Conversations are nil")
	assert.Equal(suite.T(), 0, len(conversations), "Length of conversations is not 0")
}

func (suite *MessageServiceUnitTestSuite) TestMessageService_GetConversationsByUser_ReturnsListOfConversations() {
	conversations := suite.service.GetConversationsByUser(2, context.TODO())

	assert.NotNil(suite.T(), conversations, "Conversations are nil")
	assert.Equal(suite.T(), 1, len(conversations), "Length of conversations is not 1")
}

func (suite *MessageServiceUnitTestSuite) TestMessageService_GetConversationByUsers_ReturnsError() {
	id := uint(0)

	conversation, err := suite.service.GetConversationByUsers(1, 2, context.TODO())

	assert.Equal(suite.T(), conversation.User1, id, "User1 is not 0")
	assert.Equal(suite.T(), conversation.User2, id, "User2 is not 0")
	assert.NotNil(suite.T(), err, "Error are nil")
}

func (suite *MessageServiceUnitTestSuite) TestMessageService_GetConversationByUsers_ReturnsConversation() {
	user1 := uint(2)
	user2 := uint(3)

	conversation, err := suite.service.GetConversationByUsers(2, 3, context.TODO())

	assert.Equal(suite.T(), conversation.User1, user1, "User1 is not 2")
	assert.Equal(suite.T(), conversation.User2, user2, "User2 is not 3")
	assert.Nil(suite.T(), err, "Error are not nil")
}

func (suite *MessageServiceUnitTestSuite) TestMessageService_CreateNewMessage_ReturnsMessage_And_CreatesNewConversation() {
	user1 := uint(1)
	user2 := uint(2)

	messageDto := dto.MessageDto{
		From: user1,
		To:   user2,
		Text: "Text",
	}

	message, err := suite.service.CreateNewMessage(messageDto, context.TODO())

	assert.Nil(suite.T(), err, "Error is not nil")
	assert.NotNil(suite.T(), message, "Conversation is nil")
	assert.Equal(suite.T(), user1, message.From, "User1 is not as expected")
	assert.Equal(suite.T(), user2, message.To, "User2 is not as expected")
	assert.Equal(suite.T(), "Text", message.Text, "Text is not as expected")
	assert.NotNil(suite.T(), message.ConversationId, "ConversationId is nil")
	assert.NotNil(suite.T(), message.CreatedAt, "CreatedAt is nil")
	assert.NotNil(suite.T(), message.UpdatedAt, "UpdatedAt is nil")
	assert.NotNil(suite.T(), message.ID, "ID is nil")
}

func (suite *MessageServiceUnitTestSuite) TestMessageService_CreateNewMessage_ReturnsMessage_And_DoesNotCreateNewConversation() {
	user1 := uint(3)
	user2 := uint(4)

	messageDto := dto.MessageDto{
		From: user1,
		To:   user2,
		Text: "Text",
	}

	message, err := suite.service.CreateNewMessage(messageDto, context.TODO())

	assert.Nil(suite.T(), err, "Error is not nil")
	assert.NotNil(suite.T(), message, "Conversation is nil")
	assert.Equal(suite.T(), user1, message.From, "User1 is not as expected")
	assert.Equal(suite.T(), user2, message.To, "User2 is not as expected")
	assert.Equal(suite.T(), "Text", message.Text, "Text is not as expected")
	assert.NotNil(suite.T(), message.ConversationId, "ConversationId is nil")
	assert.NotNil(suite.T(), message.CreatedAt, "CreatedAt is nil")
	assert.NotNil(suite.T(), message.UpdatedAt, "UpdatedAt is nil")
	assert.NotNil(suite.T(), message.ID, "ID is nil")
}
