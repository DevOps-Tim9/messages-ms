package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Conversation struct {
	ID          primitive.ObjectID `bson:"_id"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
	User1       uint               `bson:"user1"`
	User2       uint               `bson:"user2"`
	LastMessage Message            `bson:"lastMessage, omitempty"`
}
