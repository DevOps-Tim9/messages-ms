package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	CreatedAt      time.Time          `json:"createdAt" bson:"created_at"`
	UpdatedAt      time.Time          `json:"updatedAt" bson:"updated_at"`
	From           uint               `json:"from" bson:"from"`
	To             uint               `json:"to" bson:"to"`
	ConversationId primitive.ObjectID `json:"conversationId" bson:"conversationId"`
	Text           string             `json:"text" bson:"text"`
}
