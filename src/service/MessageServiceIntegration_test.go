package service

import (
	"context"
	"fmt"
	"messages-ms/src/entity"
	"messages-ms/src/repository"
	"messages-ms/src/utils"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MessageServiceIntegrationTestSuite struct {
	suite.Suite
	service       MessageService
	db            *mongo.Database
	messages      []entity.Message
	conversations []entity.Conversation
	id            primitive.ObjectID
}

func (suite *MessageServiceIntegrationTestSuite) SetupSuite() {
	host := os.Getenv("DATABASE_DOMAIN")
	name := os.Getenv("DATABASE_SCHEMA")
	port := os.Getenv("DATABASE_PORT")

	db, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", host, port)))

	if err != nil {
		panic(err)
	}

	db.Database(name).Collection("messages")
	db.Database(name).Collection("conversations")

	database := db.Database(name)

	messageRepository := repository.MessageRepository{Database: database}
	conversationRepository := repository.ConversationRepository{Database: database}

	suite.db = database
	suite.service = MessageService{
		ConversationRepository: conversationRepository,
		MessageRepository:      messageRepository,
		Logger:                 utils.Logger(),
		RabbitMQChannel:        nil,
	}

	id := primitive.NewObjectID()

	suite.id = id
	suite.conversations = []entity.Conversation{
		{
			ID:    id,
			User1: 1,
			User2: 2,
		},
	}

	suite.messages = []entity.Message{
		{
			From:           1,
			To:             2,
			Text:           "First message",
			ConversationId: id,
		},
		{
			From:           2,
			To:             1,
			Text:           "Second message - answer",
			ConversationId: id},
	}

	suite.db.Collection("conversations").InsertOne(context.TODO(), &suite.conversations[0])
	suite.db.Collection("messages").InsertOne(context.TODO(), &suite.messages[0])
	suite.db.Collection("messages").InsertOne(context.TODO(), &suite.messages[1])
}

func TestMessageServiceIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(MessageServiceIntegrationTestSuite))
}

func (suite *MessageServiceIntegrationTestSuite) TestIntegrationMessageService_GetConversationByUsers_ConversationDoesNotExist() {
	user1 := uint(50)
	user2 := uint(60)

	_, err := suite.service.GetConversationByUsers(user1, user2, context.TODO())

	assert.NotNil(suite.T(), err)
}

func (suite *MessageServiceIntegrationTestSuite) TestIntegrationMessageService_GetConversationByUsers_ConversationDoesExist() {
	user1 := uint(1)
	user2 := uint(2)

	conversation, err := suite.service.GetConversationByUsers(user1, user2, context.TODO())

	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), conversation)
}

func (suite *MessageServiceIntegrationTestSuite) TestIntegrationMessageService_GetConversationsByUser_ConversationsExist() {
	user1 := uint(1)

	conversations := suite.service.GetConversationsByUser(user1, context.TODO())

	assert.NotNil(suite.T(), conversations)
	assert.Equal(suite.T(), len(conversations), 1)
}

func (suite *MessageServiceIntegrationTestSuite) TestIntegrationMessageService_GetConversationsByUser_ConversationsDontExist() {
	user1 := uint(9999)

	conversations := suite.service.GetConversationsByUser(user1, context.TODO())

	assert.Nil(suite.T(), conversations)
	assert.Equal(suite.T(), len(conversations), 0)
}

func (suite *MessageServiceIntegrationTestSuite) TestIntegrationMessageService_GetMesssagesByConversation_MessagesExist() {
	conversations := suite.service.GetMesssagesByConversation(suite.id.Hex(), context.TODO())

	assert.NotNil(suite.T(), conversations)
	assert.Equal(suite.T(), 1, len(conversations))
}

func (suite *MessageServiceIntegrationTestSuite) TestIntegrationMessageService_GetMesssagesByConversation_MessagesNotExist() {
	conversations := suite.service.GetMesssagesByConversation(primitive.NewObjectID().Hex(), context.TODO())

	assert.Nil(suite.T(), conversations)
	assert.Equal(suite.T(), len(conversations), 0)
}
