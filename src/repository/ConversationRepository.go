package repository

import (
	"context"
	"messages-ms/src/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IConversationRepository interface {
	Create(entity.Conversation) (entity.Conversation, error)
	Update(entity.Conversation) error
	GetConversationByUsers(uint, uint) (entity.Conversation, error)
	GetConversationsByUser(uint) []entity.Conversation
}

type ConversationRepository struct {
	Database *mongo.Database
}

func (r ConversationRepository) Create(conversation entity.Conversation) (entity.Conversation, error) {
	var newConversation = entity.Conversation{}

	result, err := r.Database.Collection("conversations").InsertOne(context.TODO(), conversation)

	if err != nil {
		return entity.Conversation{}, err
	}

	document := r.Database.Collection("conversations").FindOne(context.TODO(), bson.D{{"_id", result.InsertedID}})

	if document.Err() != nil {
		return entity.Conversation{}, document.Err()
	}

	err = document.Decode(&newConversation)

	return newConversation, err
}

func (r ConversationRepository) Update(conversation entity.Conversation) error {
	_, err := r.Database.Collection("conversations").UpdateOne(context.TODO(), bson.M{"_id": bson.M{"$eq": conversation.ID}}, bson.M{
		"$set": bson.M{
			"lastMessage": conversation.LastMessage,
			"updated_at":  conversation.UpdatedAt,
		},
	})

	return err
}

func (r ConversationRepository) GetConversationByUsers(user1 uint, user2 uint) (entity.Conversation, error) {
	var conversation = entity.Conversation{}

	filter := bson.D{
		{"$or",
			bson.A{
				bson.M{
					"user1": user1,
					"user2": user2,
				},
				bson.M{
					"user1": user2,
					"user2": user1,
				},
			},
		}}

	result := r.Database.Collection("conversations").FindOne(context.TODO(), filter)

	if result.Err() != nil {
		return conversation, result.Err()
	}

	err := result.Decode(&conversation)

	return conversation, err
}

func (r ConversationRepository) GetConversationsByUser(userId uint) []entity.Conversation {
	var conversations []entity.Conversation

	filter := bson.D{
		{"$or", bson.A{
			bson.M{"user1": userId},
			bson.M{"user2": userId},
		}},
	}

	opts := options.Find().SetSort(bson.D{{"updated_at", -1}})

	cursor, err := r.Database.Collection("conversations").Find(context.TODO(), filter, opts)

	if err != nil {
		return make([]entity.Conversation, 0)
	}

	if err = cursor.All(context.TODO(), &conversations); err != nil {
		return make([]entity.Conversation, 0)
	}

	return conversations
}
