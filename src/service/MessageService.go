package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"messages-ms/src/dto"
	"messages-ms/src/entity"
	"messages-ms/src/rabbitmq"
	"messages-ms/src/repository"
	"net/http"
	"os"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IMessageService interface {
	CreateNewMessage(dto.MessageDto, context.Context) (entity.Message, error)
	GetMesssagesByConversation(string, context.Context) []entity.Message
	GetConversationsByUser(uint, context.Context) []entity.Conversation
	GetConversationByUsers(uint, uint, context.Context) (entity.Conversation, error)
}

type MessageService struct {
	MessageRepository      repository.IMessageRepository
	ConversationRepository repository.IConversationRepository
	Logger                 *logrus.Entry
	RabbitMQChannel        *amqp.Channel
}

func (s MessageService) GetMesssagesByConversation(id string, ctx context.Context) []entity.Message {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Service - Get messages for specific conversation")

	defer span.Finish()

	s.Logger.Info("Getting all messages by conversation")

	return s.MessageRepository.GetMesssagesByConversation(id, ctx)
}

func (s MessageService) GetConversationsByUser(id uint, ctx context.Context) []entity.Conversation {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Service - Get conversations for specific user")

	defer span.Finish()

	s.Logger.Info("Getting all conversations for specified user")

	return s.ConversationRepository.GetConversationsByUser(id, ctx)
}

func (s MessageService) GetConversationByUsers(user1 uint, user2 uint, ctx context.Context) (entity.Conversation, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Service - Get conversation by users")

	defer span.Finish()

	s.Logger.Info("Getting conversation by participants")

	return s.ConversationRepository.GetConversationByUsers(user1, user2, ctx)
}

func (s MessageService) AddNotification(message dto.MessageDto, ctx context.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx, "Service - Send notification to inform of new message")

	defer span.Finish()

	userFrom := dto.UserResponseDTO{}
	userTo := dto.UserResponseDTO{}

	endpointFrom := fmt.Sprintf("http://%s/users/%d", os.Getenv("USER_SERVICE_DOMAIN"), message.From)

	reqFrom, err := http.NewRequest("GET", endpointFrom, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	resFrom, err := http.DefaultClient.Do(reqFrom)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer resFrom.Body.Close()

	if resFrom.StatusCode != 200 {
		fmt.Println(resFrom.StatusCode)
		return
	} else {
		b, _ := io.ReadAll(resFrom.Body)
		errr := json.Unmarshal(b, &userFrom)
		if errr != nil {
			fmt.Println(errr.Error())
		}
	}

	endpointTo := fmt.Sprintf("http://%s/users/%d", os.Getenv("USER_SERVICE_DOMAIN"), message.To)

	reqTo, _ := http.NewRequest("GET", endpointTo, nil)

	resTo, err := http.DefaultClient.Do(reqTo)
	if err != nil {
		return
	}
	defer resTo.Body.Close()

	if resTo.StatusCode != 200 {
		fmt.Println(resTo.StatusCode)
		return
	} else {
		b, _ := io.ReadAll(resTo.Body)
		errr := json.Unmarshal(b, &userTo)
		if errr != nil {
			fmt.Println(errr.Error())
		}
	}

	messageType := dto.Message
	notification := dto.NotificationDTO{Message: fmt.Sprintf("%s messaged you.", userFrom.Username), UserAuth0ID: userTo.Auth0ID, NotificationType: &messageType}

	rabbitmq.AddNotification(&notification, s.RabbitMQChannel)
}

func (s MessageService) CreateNewMessage(dto dto.MessageDto, ctx context.Context) (entity.Message, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Service - Create message")

	defer span.Finish()

	s.Logger.Info("Saving new message")

	conversation, err := s.ConversationRepository.GetConversationByUsers(dto.From, dto.To, ctx)

	if err != nil {
		conversation = entity.Conversation{
			ID:        primitive.NewObjectID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			User1:     dto.From,
			User2:     dto.To,
		}

		s.Logger.Info("Creating new conversation")

		newConversation, _ := s.ConversationRepository.Create(conversation, ctx)

		s.Logger.Info("Saving new message in DB")

		newMessage, _ := s.MessageRepository.Create(entity.Message{
			ID:             primitive.NewObjectID(),
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
			From:           dto.From,
			To:             dto.To,
			Text:           dto.Text,
			ConversationId: newConversation.ID,
		}, ctx)

		newConversation.LastMessage = newMessage
		newConversation.UpdatedAt = time.Now()

		s.ConversationRepository.Update(newConversation, ctx)

		s.AddNotification(dto, ctx)

		return newMessage, nil
	} else {
		s.Logger.Info("Saving new message in DB")

		newMessage, err := s.MessageRepository.Create(entity.Message{
			ID:             primitive.NewObjectID(),
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
			From:           dto.From,
			To:             dto.To,
			Text:           dto.Text,
			ConversationId: conversation.ID,
		}, ctx)

		conversation.LastMessage = newMessage
		conversation.UpdatedAt = time.Now()

		s.ConversationRepository.Update(conversation, ctx)

		s.AddNotification(dto, ctx)

		return newMessage, err
	}
}
