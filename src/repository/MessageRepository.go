package repository

import (
	"context"
	"messages-ms/src/entity"

	"github.com/opentracing/opentracing-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IMessageRepository interface {
	Create(entity.Message, context.Context) (entity.Message, error)
	Update(entity.Message, context.Context) error
	GetMesssagesByConversation(string, context.Context) []entity.Message
}

type MessageRepository struct {
	Database *mongo.Database
}

func (r MessageRepository) Create(message entity.Message, ctx context.Context) (entity.Message, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "Repository - Crete new message")

	defer span.Finish()

	var newMessage = entity.Message{}

	result, err := r.Database.Collection("messages").InsertOne(context.TODO(), message)

	if err != nil {
		return entity.Message{}, err
	}

	document := r.Database.Collection("messages").FindOne(context.TODO(), bson.D{{"_id", result.InsertedID}})

	if document.Err() != nil {
		return entity.Message{}, document.Err()
	}

	err = document.Decode(&newMessage)

	return newMessage, err
}

func (r MessageRepository) Update(message entity.Message, ctx context.Context) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "Repository - Update message")

	defer span.Finish()

	_, err := r.Database.Collection("messages").UpdateOne(context.TODO(), bson.D{{"_id", message.ID}}, message)

	return err
}

func (r MessageRepository) GetMesssagesByConversation(conversationId string, ctx context.Context) []entity.Message {
	span, _ := opentracing.StartSpanFromContext(ctx, "Repository - Get messages for specific conversation")

	defer span.Finish()

	var messages []entity.Message

	opts := options.Find().SetSort(bson.D{{"created_at", -1}})

	id, _ := primitive.ObjectIDFromHex(conversationId)

	cursor, err := r.Database.Collection("messages").Find(context.TODO(), bson.D{{Key: "conversationId", Value: id}}, opts)

	if err != nil {
		return make([]entity.Message, 0)
	}

	if err = cursor.All(context.TODO(), &messages); err != nil {
		return make([]entity.Message, 0)
	}

	return messages
}
